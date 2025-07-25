/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: BUSL-1.1
*/

package deploy

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImageVersion(t *testing.T) {
	testCases := map[string]struct {
		imageReference string
		createFile     [2]string
		wantVersion    string
		wantErr        bool
	}{
		"version found in /etc": {
			imageReference: "some-reference",
			createFile:     [2]string{"/host/etc/os-release", osRelease},
			wantVersion:    "v2.3.0",
		},
		"version found in /usr/lib": {
			imageReference: "some-reference",
			createFile:     [2]string{"/host/usr/lib/os-release", osRelease},
			wantVersion:    "v2.3.0",
		},
		"version not found": {
			imageReference: "some-reference",
			wantErr:        true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			fs := afero.NewMemMapFs()
			if tc.createFile[0] != "" {
				err := afero.WriteFile(fs, tc.createFile[0], []byte(tc.createFile[1]), os.ModePerm)
				require.NoError(err)
			}

			imageInfo := &ImageInfo{
				fs: &afero.Afero{Fs: fs},
			}

			version, err := imageInfo.ImageVersion()
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			assert.Equal(tc.wantVersion, version)
		})
	}
}

func TestGetOSReleaseImageVersion(t *testing.T) {
	testCases := map[string]struct {
		path        string
		wantVersion string
		wantErr     bool
	}{
		"version found": {
			path:        "os-release",
			wantVersion: "v2.3.0",
		},
		"invalid path": {
			path:    "not/a/real/path",
			wantErr: true,
		},
		"empty file": {
			path:    "empty",
			wantErr: true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			fs := afero.NewMemMapFs()
			err := afero.WriteFile(fs, "os-release", []byte(osRelease), os.ModePerm)
			require.NoError(err)
			err = afero.WriteFile(fs, "empty", []byte{}, os.ModePerm)
			require.NoError(err)

			imageInfo := &ImageInfo{
				fs: &afero.Afero{Fs: fs},
			}

			version, err := imageInfo.getOSReleaseImageVersion(tc.path)
			if tc.wantErr {
				assert.Error(err)
				return
			}
			require.NoError(err)
			assert.Equal(tc.wantVersion, version)
		})
	}
}

func TestParseOSRelease(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)
	osReleaseMap, err := parseOSRelease(bufio.NewScanner(strings.NewReader(osRelease)))
	require.NoError(err)
	assert.Equal(wantMap, osReleaseMap)
}

const osRelease = `
# Some comment
# Some empty lines below


SINGLE_QUOTED_VALUE='WOW! This is a single quoted value!'
DOUBLE_QUOTED_VALUE="WOW! This is a double quoted value!"
ESCAPED_BACKSLASH='This is a string with an escaped backslash: \\'
ESCAPED_DOLLAR='This is a string with an escaped dollar: \$'
ESCAPED_DOUBLE_QUOTE='This is a string with an escaped double quote: \"'
ESCAPED_SINGLE_QUOTE="This is a string with an escaped single quote: \'"
NAME="Fedora Linux"
VERSION="37 (Thirty Seven)"
ID=fedora
PRETTY_NAME="Fedora Linux 37 (Thirty Seven)"
ANSI_COLOR="0;38;2;60;110;180"
VERSION_ID=37
VERSION_CODENAME=""
PLATFORM_ID="platform:f37"
LOGO=fedora-logo-icon
CPE_NAME="cpe:/o:fedoraproject:fedora:37"
DEFAULT_HOSTNAME="fedora"
HOME_URL="https://fedoraproject.org/"
DOCUMENTATION_URL="https://docs.fedoraproject.org/en-US/fedora/f37/system-administrators-guide/"
SUPPORT_URL="https://ask.fedoraproject.org/"
BUG_REPORT_URL="https://bugzilla.redhat.com/"
REDHAT_BUGZILLA_PRODUCT="Fedora"
REDHAT_BUGZILLA_PRODUCT_VERSION=37
REDHAT_SUPPORT_PRODUCT="Fedora"
REDHAT_SUPPORT_PRODUCT_VERSION=37
IMAGE_ID="constellation"
IMAGE_VERSION="v2.3.0"
`

var wantMap = map[string]string{
	"NAME":                            `Fedora Linux`,
	"VERSION":                         `37 (Thirty Seven)`,
	"ID":                              `fedora`,
	"SINGLE_QUOTED_VALUE":             `WOW! This is a single quoted value!`,
	"DOUBLE_QUOTED_VALUE":             `WOW! This is a double quoted value!`,
	"ESCAPED_BACKSLASH":               `This is a string with an escaped backslash: \`,
	"ESCAPED_DOLLAR":                  `This is a string with an escaped dollar: $`,
	"ESCAPED_DOUBLE_QUOTE":            `This is a string with an escaped double quote: "`,
	"ESCAPED_SINGLE_QUOTE":            `This is a string with an escaped single quote: '`,
	"VERSION_ID":                      `37`,
	"VERSION_CODENAME":                ``,
	"PLATFORM_ID":                     `platform:f37`,
	"PRETTY_NAME":                     `Fedora Linux 37 (Thirty Seven)`,
	"ANSI_COLOR":                      `0;38;2;60;110;180`,
	"LOGO":                            `fedora-logo-icon`,
	"CPE_NAME":                        `cpe:/o:fedoraproject:fedora:37`,
	"DEFAULT_HOSTNAME":                `fedora`,
	"HOME_URL":                        `https://fedoraproject.org/`,
	"DOCUMENTATION_URL":               `https://docs.fedoraproject.org/en-US/fedora/f37/system-administrators-guide/`,
	"SUPPORT_URL":                     `https://ask.fedoraproject.org/`,
	"BUG_REPORT_URL":                  `https://bugzilla.redhat.com/`,
	"REDHAT_BUGZILLA_PRODUCT":         `Fedora`,
	"REDHAT_BUGZILLA_PRODUCT_VERSION": `37`,
	"REDHAT_SUPPORT_PRODUCT":          `Fedora`,
	"REDHAT_SUPPORT_PRODUCT_VERSION":  `37`,
	"IMAGE_ID":                        `constellation`,
	"IMAGE_VERSION":                   `v2.3.0`,
}
