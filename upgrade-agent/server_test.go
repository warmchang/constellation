/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package upgradeagent

import (
	"testing"

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
