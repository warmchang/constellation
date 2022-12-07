/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package upgradeagent

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"

	"github.com/edgelesssys/constellation/v2/internal/constants"
	"github.com/edgelesssys/constellation/v2/internal/file"
	"github.com/edgelesssys/constellation/v2/internal/installer"
	"github.com/edgelesssys/constellation/v2/internal/logger"
	"github.com/edgelesssys/constellation/v2/internal/versions"
	"github.com/edgelesssys/constellation/v2/upgrade-agent/upgradeproto"
	"google.golang.org/grpc"
)

var versionRegexp = regexp.MustCompile(`^v\d{1}\.\d{1,2}\.\d{1,2}$`)

// Server is the upgrade-agent server.
type Server struct {
	file       file.Handler
	grpcServer serveStopper
	log        *logger.Logger
	upgradeproto.UnimplementedUpdateServer
}

// New creates a new upgrade-agent server.
func New(log *logger.Logger, fileHandler file.Handler) (*Server, error) {
	log = log.Named("upgradeServer")

	server := &Server{
		log:  log,
		file: fileHandler,
	}

	grpcServer := grpc.NewServer(
		log.Named("gRPC").GetServerUnaryInterceptor(),
	)
	upgradeproto.RegisterUpdateServer(grpcServer, server)

	server.grpcServer = grpcServer
	return server, nil
}

// Run starts the upgrade-agent server on the given port, using the provided protocol and socket address.
func (s *Server) Run(protocol string, sockAddr string) error {
	grpcServer := grpc.NewServer()

	upgradeproto.RegisterUpdateServer(grpcServer, s)

	cleanup := func() error {
		if _, err := os.Stat(sockAddr); err == nil {
			if err := os.RemoveAll(sockAddr); err != nil {
				return err
			}
		}
		return nil
	}
	if err := cleanup(); err != nil {
		return fmt.Errorf("failed to clean socket file: %s", err)
	}

	lis, err := net.Listen(protocol, sockAddr)
	if err != nil {
		return fmt.Errorf("failed to listen: %s", err)
	}

	s.log.Infof("Starting")
	return grpcServer.Serve(lis)
}

// Stop stops the upgrade-agent server gracefully.
func (s *Server) Stop() {
	s.log.Infof("Stopping")

	s.grpcServer.GracefulStop()

	s.log.Infof("Stopped")
}

// ExecuteUpdate installs & verifies the provided kubeadm, then executes `kubeadm upgrade plan` & `kubeadm upgrade apply {wanted_Kubernetes_Version}` to upgrade to the specified version.
func (s *Server) ExecuteUpdate(ctx context.Context, updateRequest *upgradeproto.ExecuteUpdateRequest) (*upgradeproto.ExecuteUpdateResponse, error) {
	s.log.Infof("Upgrade to Kubernetes version started: %s", updateRequest.WantedKubernetesVersion)
	installer := installer.NewOSInstaller()
	if err := prepareUpdate(ctx, installer, updateRequest); err != nil {
		return &upgradeproto.ExecuteUpdateResponse{}, err
	}

	upgradeCmd := exec.CommandContext(ctx, "kubeadm", "upgrade", "plan")
	if err := upgradeCmd.Run(); err != nil {
		return &upgradeproto.ExecuteUpdateResponse{}, err
	}

	applyCmd := exec.CommandContext(ctx, "kubeadm", "upgrade", "apply", updateRequest.WantedKubernetesVersion)
	if err := applyCmd.Run(); err != nil {
		return &upgradeproto.ExecuteUpdateResponse{}, err
	}

	s.log.Infof("Upgrade to Kubernetes version succeeded: %s", updateRequest.WantedKubernetesVersion)
	return &upgradeproto.ExecuteUpdateResponse{}, nil
}

// prepareUpdate downloads & installs the specified kubeadm version and verifies the desired Kubernetes version.
func prepareUpdate(ctx context.Context, installer osInstaller, updateRequest *upgradeproto.ExecuteUpdateRequest) error {
	// download & install the kubeadm binary
	err := installer.Install(ctx, versions.ComponentVersion{
		URL:         updateRequest.KubeadmUrl,
		Hash:        updateRequest.KubeadmHash,
		InstallPath: constants.KubeadmPath,
		Extract:     false,
	})
	if err != nil {
		return err
	}

	err = verifyVersion(updateRequest.WantedKubernetesVersion)
	if err != nil {
		return err
	}

	return nil
}

// verifyVersion verifies the provided Kubernetes version.
func verifyVersion(version string) error {
	if !versionRegexp.MatchString(version) {
		return fmt.Errorf("invalid kubernetes version: %s", version)
	}
	return nil
}

type osInstaller interface {
	// Install downloads, installs and verifies the kubernetes component.
	Install(ctx context.Context, kubernetesComponent versions.ComponentVersion) error
}

type serveStopper interface {
	// Serve starts the server.
	Serve(lis net.Listener) error
	// GracefulStop stops the server and blocks until all requests are done.
	GracefulStop()
}
