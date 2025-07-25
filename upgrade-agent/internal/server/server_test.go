/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: BUSL-1.1
*/

package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/edgelesssys/constellation/v2/internal/versions/components"
	"github.com/edgelesssys/constellation/v2/upgrade-agent/upgradeproto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVersionVerifier(t *testing.T) {
	testCases := map[string]struct {
		versionString string
		wantErr       bool
	}{
		"valid version": {
			versionString: "v1.1.1",
		},
		"v prefix missing": {
			versionString: "1.1.1",
			wantErr:       true,
		},
		"invalid space": {
			versionString: "v 1.1.1",
			wantErr:       true,
		},
		"invalid version": {
			versionString: "v1.1.1a",
			wantErr:       true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			err := verifyVersion(tc.versionString)
			if tc.wantErr {
				assert.Error(err)
				return
			}

			require.NoError(err)
		})
	}
}

func TestPrepareUpdate(t *testing.T) {
	invalidUpgradeRequest := &upgradeproto.ExecuteUpdateRequest{
		WantedKubernetesVersion: "1337",
	}
	slimUpdateRequest := &upgradeproto.ExecuteUpdateRequest{
		WantedKubernetesVersion: "v1.1.1",
	}
	updateRequest := &upgradeproto.ExecuteUpdateRequest{
		WantedKubernetesVersion: "v1.1.1",
		KubernetesComponents: []*components.Component{
			{
				Url:         "http://example.com/kubeadm",
				Hash:        "sha256:foo",
				InstallPath: "/tmp/kubeadm",
			},
		},
	}
	combinedStyleUpdateRequest := &upgradeproto.ExecuteUpdateRequest{
		WantedKubernetesVersion: "v1.1.1",
		KubernetesComponents: []*components.Component{
			{
				Url:         "data:application/octet-stream,foo",
				InstallPath: "/tmp/foo",
			},
		},
	}
	testCases := map[string]struct {
		installer     osInstaller
		updateRequest *upgradeproto.ExecuteUpdateRequest
		wantErr       bool
	}{
		"works": {
			installer:     stubOsInstaller{},
			updateRequest: slimUpdateRequest,
		},
		"invalid version string": {
			installer:     stubOsInstaller{},
			updateRequest: invalidUpgradeRequest,
			wantErr:       true,
		},
		"install error": {
			installer:     stubOsInstaller{InstallErr: fmt.Errorf("install error")},
			updateRequest: updateRequest,
			wantErr:       true,
		},
		"new style works": {
			installer:     stubOsInstaller{},
			updateRequest: updateRequest,
		},
		"new style install error": {
			installer:     stubOsInstaller{InstallErr: fmt.Errorf("install error")},
			updateRequest: updateRequest,
			wantErr:       true,
		},
		"combined style works": {
			installer:     stubOsInstaller{},
			updateRequest: combinedStyleUpdateRequest,
		},
		"combined style install error": {
			installer:     stubOsInstaller{InstallErr: fmt.Errorf("install error")},
			updateRequest: combinedStyleUpdateRequest,
			wantErr:       true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			err := prepareUpdate(t.Context(), tc.installer, tc.updateRequest)
			if tc.wantErr {
				assert.Error(err)
				return
			}

			require.NoError(err)
		})
	}
}

type stubOsInstaller struct {
	InstallErr error
}

func (s stubOsInstaller) Install(_ context.Context, _ *components.Component) error {
	return s.InstallErr
}
