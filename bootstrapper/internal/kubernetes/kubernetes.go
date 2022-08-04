package kubernetes

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"

	"github.com/edgelesssys/constellation/bootstrapper/internal/kubernetes/k8sapi"
	"github.com/edgelesssys/constellation/bootstrapper/internal/kubernetes/k8sapi/resources"
	"github.com/edgelesssys/constellation/bootstrapper/role"
	"github.com/edgelesssys/constellation/bootstrapper/util"
	"github.com/edgelesssys/constellation/internal/cloud/metadata"
	"github.com/edgelesssys/constellation/internal/constants"
	"github.com/edgelesssys/constellation/internal/logger"
	"github.com/edgelesssys/constellation/internal/versions"
	"github.com/spf13/afero"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubeadm "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta3"
)

// configReader provides kubeconfig as []byte.
type configReader interface {
	ReadKubeconfig() ([]byte, error)
}

// configurationProvider provides kubeadm init and join configuration.
type configurationProvider interface {
	InitConfiguration(externalCloudProvider bool, k8sVersion versions.ValidK8sVersion) k8sapi.KubeadmInitYAML
	JoinConfiguration(externalCloudProvider bool) k8sapi.KubeadmJoinYAML
}

// KubeWrapper implements Cluster interface.
type KubeWrapper struct {
	cloudProvider           string
	clusterUtil             clusterUtil
	configProvider          configurationProvider
	client                  k8sapi.Client
	kubeconfigReader        configReader
	cloudControllerManager  CloudControllerManager
	cloudNodeManager        CloudNodeManager
	clusterAutoscaler       ClusterAutoscaler
	providerMetadata        ProviderMetadata
	initialMeasurementsJSON []byte
	getIPAddr               func() (string, error)
}

// New creates a new KubeWrapper with real values.
func New(cloudProvider string, clusterUtil clusterUtil, configProvider configurationProvider, client k8sapi.Client, cloudControllerManager CloudControllerManager,
	cloudNodeManager CloudNodeManager, clusterAutoscaler ClusterAutoscaler, providerMetadata ProviderMetadata, initialMeasurementsJSON []byte,
) *KubeWrapper {
	return &KubeWrapper{
		cloudProvider:           cloudProvider,
		clusterUtil:             clusterUtil,
		configProvider:          configProvider,
		client:                  client,
		kubeconfigReader:        &KubeconfigReader{fs: afero.Afero{Fs: afero.NewOsFs()}},
		cloudControllerManager:  cloudControllerManager,
		cloudNodeManager:        cloudNodeManager,
		clusterAutoscaler:       clusterAutoscaler,
		providerMetadata:        providerMetadata,
		initialMeasurementsJSON: initialMeasurementsJSON,
		getIPAddr:               util.GetIPAddr,
	}
}

// InitCluster initializes a new Kubernetes cluster and applies pod network provider.
func (k *KubeWrapper) InitCluster(
	ctx context.Context, autoscalingNodeGroups []string, cloudServiceAccountURI, versionString string,
	measurementSalt []byte, kmsConfig resources.KMSConfig, sshUsers map[string]string, log *logger.Logger,
) ([]byte, error) {
	k8sVersion, err := versions.NewValidK8sVersion(versionString)
	if err != nil {
		return nil, err
	}
	log.With(zap.String("version", string(k8sVersion))).Infof("Installing Kubernetes components")
	if err := k.clusterUtil.InstallComponents(ctx, k8sVersion); err != nil {
		return nil, err
	}

	ip, err := k.getIPAddr()
	if err != nil {
		return nil, err
	}
	nodeName := ip
	var providerID string
	var instance metadata.InstanceMetadata
	var publicIP string
	var nodePodCIDR string
	var subnetworkPodCIDR string
	var controlPlaneEndpointIP string // this is the IP in "kubeadm init --control-plane-endpoint=<IP/DNS>:<port>" hence the unfortunate name
	var nodeIP string
	var validIPs []net.IP

	// Step 1: retrieve cloud metadata for Kubernetes configuration
	if k.providerMetadata.Supported() {
		log.Infof("Retrieving node metadata")
		instance, err = k.providerMetadata.Self(ctx)
		if err != nil {
			return nil, fmt.Errorf("retrieving own instance metadata failed: %w", err)
		}
		if instance.VPCIP != "" {
			validIPs = append(validIPs, net.ParseIP(instance.VPCIP))
		}
		if instance.PublicIP != "" {
			validIPs = append(validIPs, net.ParseIP(instance.PublicIP))
		}
		nodeName = k8sCompliantHostname(instance.Name)
		providerID = instance.ProviderID
		nodeIP = instance.VPCIP
		publicIP = instance.PublicIP

		if len(instance.AliasIPRanges) > 0 {
			nodePodCIDR = instance.AliasIPRanges[0]
		}
		subnetworkPodCIDR, err = k.providerMetadata.GetSubnetworkCIDR(ctx)
		if err != nil {
			return nil, fmt.Errorf("retrieving subnetwork CIDR failed: %w", err)
		}
		controlPlaneEndpointIP = publicIP
		if k.providerMetadata.SupportsLoadBalancer() {
			controlPlaneEndpointIP, err = k.providerMetadata.GetLoadBalancerIP(ctx)
			if err != nil {
				return nil, fmt.Errorf("retrieving load balancer IP failed: %w", err)
			}
			if k.cloudProvider == "gcp" {
				if err := manuallySetLoadbalancerIP(ctx, controlPlaneEndpointIP); err != nil {
					return nil, fmt.Errorf("setting load balancer IP failed: %w", err)
				}
			}
		}
	}
	log.With(
		zap.String("nodeName", nodeName),
		zap.String("providerID", providerID),
		zap.String("nodeIP", nodeIP),
		zap.String("controlPlaneEndpointIP", controlPlaneEndpointIP),
		zap.String("podCIDR", subnetworkPodCIDR),
	).Infof("Setting information for node")

	// Step 2: configure kubeadm init config
	initConfig := k.configProvider.InitConfiguration(k.cloudControllerManager.Supported(), k8sVersion)
	initConfig.SetNodeIP(nodeIP)
	initConfig.SetCertSANs([]string{publicIP, nodeIP})
	initConfig.SetNodeName(nodeName)
	initConfig.SetProviderID(providerID)
	initConfig.SetControlPlaneEndpoint(controlPlaneEndpointIP)
	initConfigYAML, err := initConfig.Marshal()
	if err != nil {
		return nil, fmt.Errorf("encoding kubeadm init configuration as YAML: %w", err)
	}
	log.Infof("Initializing Kubernetes cluster")
	if err := k.clusterUtil.InitCluster(ctx, initConfigYAML, nodeName, validIPs, log); err != nil {
		return nil, fmt.Errorf("kubeadm init: %w", err)
	}
	kubeConfig, err := k.GetKubeconfig()
	if err != nil {
		return nil, fmt.Errorf("reading kubeconfig after cluster initialization: %w", err)
	}
	k.client.SetKubeconfig(kubeConfig)

	// Step 3: configure & start kubernetes controllers
	log.Infof("Starting Kubernetes controllers and deployments")
	setupPodNetworkInput := k8sapi.SetupPodNetworkInput{
		CloudProvider:     k.cloudProvider,
		NodeName:          nodeName,
		FirstNodePodCIDR:  nodePodCIDR,
		SubnetworkPodCIDR: subnetworkPodCIDR,
		ProviderID:        providerID,
	}
	if err = k.clusterUtil.SetupPodNetwork(ctx, setupPodNetworkInput, k.client); err != nil {
		return nil, fmt.Errorf("setting up pod network: %w", err)
	}

	kms := resources.NewKMSDeployment(k.cloudProvider, kmsConfig)
	if err = k.clusterUtil.SetupKMS(k.client, kms); err != nil {
		return nil, fmt.Errorf("setting up kms: %w", err)
	}

	if err := k.setupJoinService(k.cloudProvider, k.initialMeasurementsJSON, measurementSalt); err != nil {
		return nil, fmt.Errorf("setting up join service failed: %w", err)
	}

	if err := k.setupCCM(ctx, subnetworkPodCIDR, cloudServiceAccountURI, instance, k8sVersion); err != nil {
		return nil, fmt.Errorf("setting up cloud controller manager: %w", err)
	}
	if err := k.setupCloudNodeManager(k8sVersion); err != nil {
		return nil, fmt.Errorf("setting up cloud node manager: %w", err)
	}

	if err := k.setupClusterAutoscaler(instance, cloudServiceAccountURI, autoscalingNodeGroups, k8sVersion); err != nil {
		return nil, fmt.Errorf("setting up cluster autoscaler: %w", err)
	}

	accessManager := resources.NewAccessManagerDeployment(sshUsers)
	if err := k.clusterUtil.SetupAccessManager(k.client, accessManager); err != nil {
		return nil, fmt.Errorf("failed to setup access-manager: %w", err)
	}

	if err := k.clusterUtil.SetupVerificationService(
		k.client, resources.NewVerificationDaemonSet(k.cloudProvider),
	); err != nil {
		return nil, fmt.Errorf("failed to setup verification service: %w", err)
	}

	// TODO: enable operator deployment on kubernetes 1.24 once https://github.com/medik8s/node-maintenance-operator/issues/49 is fixed
	if k8sVersion != versions.V1_24 {
		if err := k.setupOperators(ctx); err != nil {
			return nil, fmt.Errorf("setting up operators: %w", err)
		}
	}

	if k.cloudProvider == "gcp" {
		if err := k.clusterUtil.SetupGCPGuestAgent(k.client, resources.NewGCPGuestAgentDaemonset()); err != nil {
			return nil, fmt.Errorf("failed to setup gcp guest agent: %w", err)
		}
	}

	// Store the received k8sVersion in a ConfigMap, overwriting existing values (there shouldn't be any).
	// Joining nodes determine the kubernetes version they will install based on this ConfigMap.
	if err := k.setupK8sVersionConfigMap(ctx, k8sVersion); err != nil {
		return nil, fmt.Errorf("failed to setup k8s version ConfigMap: %v", err)
	}

	k.clusterUtil.FixCilium(nodeName, log)

	return k.GetKubeconfig()
}

// JoinCluster joins existing Kubernetes cluster.
func (k *KubeWrapper) JoinCluster(ctx context.Context, args *kubeadm.BootstrapTokenDiscovery, peerRole role.Role, versionString string, log *logger.Logger) error {
	k8sVersion, err := versions.NewValidK8sVersion(versionString)
	if err != nil {
		return err
	}
	log.With(zap.String("version", string(k8sVersion))).Infof("Installing Kubernetes components")
	if err := k.clusterUtil.InstallComponents(ctx, k8sVersion); err != nil {
		return err
	}

	// Step 1: retrieve cloud metadata for Kubernetes configuration
	nodeInternalIP, err := k.getIPAddr()
	if err != nil {
		return err
	}
	nodeName := nodeInternalIP
	var providerID string
	if k.providerMetadata.Supported() {
		log.Infof("Retrieving node metadata")
		instance, err := k.providerMetadata.Self(ctx)
		if err != nil {
			return fmt.Errorf("retrieving own instance metadata failed: %w", err)
		}
		providerID = instance.ProviderID
		nodeName = instance.Name
		nodeInternalIP = instance.VPCIP
	}
	nodeName = k8sCompliantHostname(nodeName)

	log.With(
		zap.String("nodeName", nodeName),
		zap.String("providerID", providerID),
		zap.String("nodeIP", nodeInternalIP),
	).Infof("Setting information for node")

	// Step 2: configure kubeadm join config

	joinConfig := k.configProvider.JoinConfiguration(k.cloudControllerManager.Supported())
	joinConfig.SetAPIServerEndpoint(args.APIServerEndpoint)
	joinConfig.SetToken(args.Token)
	joinConfig.AppendDiscoveryTokenCaCertHash(args.CACertHashes[0])
	joinConfig.SetNodeIP(nodeInternalIP)
	joinConfig.SetNodeName(nodeName)
	joinConfig.SetProviderID(providerID)
	if peerRole == role.ControlPlane {
		joinConfig.SetControlPlane(nodeInternalIP)
	}
	joinConfigYAML, err := joinConfig.Marshal()
	if err != nil {
		return fmt.Errorf("encoding kubeadm join configuration as YAML: %w", err)
	}
	log.With(zap.String("apiServerEndpoint", args.APIServerEndpoint)).Infof("Joining Kubernetes cluster")
	if err := k.clusterUtil.JoinCluster(ctx, joinConfigYAML, log); err != nil {
		return fmt.Errorf("joining cluster: %v; %w ", string(joinConfigYAML), err)
	}

	k.clusterUtil.FixCilium(nodeName, log)

	return nil
}

// GetKubeconfig returns the current nodes kubeconfig of stored on disk.
func (k *KubeWrapper) GetKubeconfig() ([]byte, error) {
	return k.kubeconfigReader.ReadKubeconfig()
}

func (k *KubeWrapper) setupJoinService(csp string, measurementsJSON, measurementSalt []byte) error {
	joinConfiguration := resources.NewJoinServiceDaemonset(csp, string(measurementsJSON), measurementSalt)

	return k.clusterUtil.SetupJoinService(k.client, joinConfiguration)
}

func (k *KubeWrapper) setupCCM(ctx context.Context, subnetworkPodCIDR, cloudServiceAccountURI string, instance metadata.InstanceMetadata, k8sVersion versions.ValidK8sVersion) error {
	if !k.cloudControllerManager.Supported() {
		return nil
	}
	ccmConfigMaps, err := k.cloudControllerManager.ConfigMaps(instance)
	if err != nil {
		return fmt.Errorf("defining ConfigMaps for CCM failed: %w", err)
	}
	ccmSecrets, err := k.cloudControllerManager.Secrets(ctx, instance.ProviderID, cloudServiceAccountURI)
	if err != nil {
		return fmt.Errorf("defining Secrets for CCM failed: %w", err)
	}
	ccmImage, err := k.cloudControllerManager.Image(k8sVersion)
	if err != nil {
		return fmt.Errorf("defining Image for CCM failed: %w", err)
	}

	cloudControllerManagerConfiguration := resources.NewDefaultCloudControllerManagerDeployment(
		k.cloudControllerManager.Name(), ccmImage, k.cloudControllerManager.Path(), subnetworkPodCIDR,
		k.cloudControllerManager.ExtraArgs(), k.cloudControllerManager.Volumes(), k.cloudControllerManager.VolumeMounts(), k.cloudControllerManager.Env(),
	)
	if err := k.clusterUtil.SetupCloudControllerManager(k.client, cloudControllerManagerConfiguration, ccmConfigMaps, ccmSecrets); err != nil {
		return fmt.Errorf("failed to setup cloud-controller-manager: %w", err)
	}

	return nil
}

func (k *KubeWrapper) setupCloudNodeManager(k8sVersion versions.ValidK8sVersion) error {
	if !k.cloudNodeManager.Supported() {
		return nil
	}
	nodeManagerImage, err := k.cloudNodeManager.Image(k8sVersion)
	if err != nil {
		return fmt.Errorf("defining Image for Node Manager failed: %w", err)
	}

	cloudNodeManagerConfiguration := resources.NewDefaultCloudNodeManagerDeployment(
		nodeManagerImage, k.cloudNodeManager.Path(), k.cloudNodeManager.ExtraArgs(),
	)
	if err := k.clusterUtil.SetupCloudNodeManager(k.client, cloudNodeManagerConfiguration); err != nil {
		return fmt.Errorf("failed to setup cloud-node-manager: %w", err)
	}

	return nil
}

func (k *KubeWrapper) setupClusterAutoscaler(instance metadata.InstanceMetadata, cloudServiceAccountURI string, autoscalingNodeGroups []string, k8sVersion versions.ValidK8sVersion) error {
	if !k.clusterAutoscaler.Supported() {
		return nil
	}
	caSecrets, err := k.clusterAutoscaler.Secrets(instance.ProviderID, cloudServiceAccountURI)
	if err != nil {
		return fmt.Errorf("defining Secrets for cluster-autoscaler failed: %w", err)
	}

	clusterAutoscalerConfiguration := resources.NewDefaultAutoscalerDeployment(k.clusterAutoscaler.Volumes(), k.clusterAutoscaler.VolumeMounts(), k.clusterAutoscaler.Env(), k8sVersion)
	clusterAutoscalerConfiguration.SetAutoscalerCommand(k.clusterAutoscaler.Name(), autoscalingNodeGroups)
	if err := k.clusterUtil.SetupAutoscaling(k.client, clusterAutoscalerConfiguration, caSecrets); err != nil {
		return fmt.Errorf("failed to setup cluster-autoscaler: %w", err)
	}

	return nil
}

// setupK8sVersionConfigMap applies a ConfigMap (cf. server-side apply) to consistently store the installed k8s version.
func (k *KubeWrapper) setupK8sVersionConfigMap(ctx context.Context, k8sVersion versions.ValidK8sVersion) error {
	config := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "k8s-version",
			Namespace: "kube-system",
		},
		Data: map[string]string{
			constants.K8sVersion: string(k8sVersion),
		},
	}

	// We do not use the client's Apply method here since we are handling a kubernetes-native type.
	// These types don't implement our custom Marshaler interface.
	if err := k.client.CreateConfigMap(ctx, config); err != nil {
		return fmt.Errorf("apply in KubeWrapper.setupK8sVersionConfigMap(..) failed with: %v", err)
	}

	return nil
}

// setupOperators deploys the operator lifecycle manager and subscriptions to operators.
func (k *KubeWrapper) setupOperators(ctx context.Context) error {
	if err := k.clusterUtil.SetupOperatorLifecycleManager(ctx, k.client, &resources.OperatorLifecycleManagerCRDs{}, &resources.OperatorLifecycleManager{}, resources.OLMCRDNames); err != nil {
		return fmt.Errorf("setting up OLM: %w", err)
	}

	if err := k.clusterUtil.SetupNodeMaintenanceOperator(k.client, resources.NewNodeMaintenanceOperatorDeployment()); err != nil {
		return fmt.Errorf("setting up node maintenance operator: %w", err)
	}

	uid, err := k.providerMetadata.UID(ctx)
	if err != nil {
		return fmt.Errorf("retrieving constellation UID: %w", err)
	}

	if err := k.clusterUtil.SetupNodeOperator(ctx, k.client, resources.NewNodeOperatorDeployment(k.cloudProvider, uid)); err != nil {
		return fmt.Errorf("setting up constellation node operator: %w", err)
	}

	return nil
}

// manuallySetLoadbalancerIP sets the loadbalancer IP of the first control plane during init.
// The GCP guest agent does this usually, but is deployed in the cluster that doesn't exist
// at this point. This is a workaround to set the loadbalancer IP manually, so kubeadm and kubelet
// can talk to the local Kubernetes API server using the loadbalancer IP.
func manuallySetLoadbalancerIP(ctx context.Context, ip string) error {
	// https://github.com/GoogleCloudPlatform/guest-agent/blob/792fce795218633bcbde505fb3457a0b24f26d37/google_guest_agent/addresses.go#L179
	if !strings.Contains(ip, "/") {
		ip = ip + "/32"
	}
	args := []string{"route", "add", "to", "local", ip, "scope", "host", "dev", "ens3", "proto", "66"}
	_, err := exec.CommandContext(ctx, "ip", args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return fmt.Errorf("ip route add (code %v) with: %s", exitErr.ExitCode(), exitErr.Stderr)
		}
		return fmt.Errorf("ip route add: %w", err)
	}
	return nil
}

// k8sCompliantHostname transforms a hostname to an RFC 1123 compliant, lowercase subdomain as required by Kubernetes node names.
// The following regex is used by k8s for validation: /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/ .
// Only a simple heuristic is used for now (to lowercase, replace underscores).
func k8sCompliantHostname(in string) string {
	hostname := strings.ToLower(in)
	hostname = strings.ReplaceAll(hostname, "_", "-")
	return hostname
}

// StartKubelet starts the kubelet service.
func (k *KubeWrapper) StartKubelet() error {
	return k.clusterUtil.StartKubelet()
}
