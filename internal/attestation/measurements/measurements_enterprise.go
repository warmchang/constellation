//go:build enterprise

/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package measurements

// Regenerate the measurements by running go generate.
// The directive can be found in the file measurements.go since it does not have
// a build tag.
// The enterprise build tag is required to validate the measurements using production
// sigstore certificates.

// revive:disable:var-naming
var (
	aws_AWSNitroTPM          = M{0: {Expected: []byte{0x73, 0x7f, 0x76, 0x7a, 0x12, 0xf5, 0x4e, 0x70, 0xee, 0xcb, 0xc8, 0x68, 0x40, 0x11, 0x32, 0x3a, 0xe2, 0xfe, 0x2d, 0xd9, 0xf9, 0x07, 0x85, 0x57, 0x79, 0x69, 0xd7, 0xa2, 0x01, 0x3e, 0x8c, 0x12}, ValidationOpt: WarnOnly}, 2: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 3: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 4: {Expected: []byte{0x94, 0x22, 0xda, 0xbc, 0xe4, 0x41, 0x60, 0x71, 0xef, 0xfe, 0xe4, 0x02, 0x67, 0x4d, 0x7c, 0x0c, 0xdd, 0x7e, 0x62, 0xb1, 0xa3, 0x26, 0xf8, 0xef, 0x8b, 0xc9, 0x00, 0x0f, 0x8b, 0xa8, 0x77, 0x49}, ValidationOpt: Enforce}, 6: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 8: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 9: {Expected: []byte{0xad, 0xca, 0x83, 0x3b, 0x2e, 0x4b, 0x6e, 0x91, 0x70, 0x89, 0xb9, 0xb5, 0xf3, 0x85, 0x08, 0x49, 0x91, 0x58, 0x6f, 0x92, 0xac, 0x2f, 0x69, 0x9c, 0x57, 0xf9, 0x0d, 0xeb, 0x88, 0xae, 0x7b, 0x85}, ValidationOpt: Enforce}, 11: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 12: {Expected: []byte{0x27, 0xdd, 0x93, 0xcd, 0xfa, 0xa5, 0xd1, 0x41, 0xc1, 0xa6, 0xb4, 0x16, 0xa1, 0xf8, 0xcc, 0xd3, 0xba, 0x0e, 0x2b, 0x33, 0xec, 0x7f, 0xbb, 0xbb, 0x8b, 0x6c, 0xca, 0xd2, 0x5e, 0x1a, 0x93, 0x6a}, ValidationOpt: Enforce}, 13: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 14: {Expected: []byte{0xd7, 0xc4, 0xcc, 0x7f, 0xf7, 0x93, 0x30, 0x22, 0xf0, 0x13, 0xe0, 0x3b, 0xde, 0xe8, 0x75, 0xb9, 0x17, 0x20, 0xb5, 0xb8, 0x6c, 0xf1, 0x75, 0x3c, 0xad, 0x83, 0x0f, 0x95, 0xe7, 0x91, 0x92, 0x6f}, ValidationOpt: WarnOnly}, 15: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}}
	aws_AWSSEVSNP            = M{0: {Expected: []byte{0x7b, 0x06, 0x8c, 0x0c, 0x3a, 0xc2, 0x9a, 0xfe, 0x26, 0x41, 0x34, 0x53, 0x6b, 0x9b, 0xe2, 0x6f, 0x1d, 0x4c, 0xcd, 0x57, 0x5b, 0x88, 0xd3, 0xc3, 0xce, 0xab, 0xf3, 0x6a, 0xc9, 0x9c, 0x02, 0x78}, ValidationOpt: WarnOnly}, 2: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 3: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 4: {Expected: []byte{0xc0, 0x6d, 0x9a, 0xe9, 0x70, 0x9f, 0x97, 0x0f, 0xf4, 0x20, 0x6e, 0x34, 0x66, 0x02, 0x9b, 0x1e, 0xfc, 0xb6, 0x2f, 0x78, 0x55, 0x20, 0x1e, 0xe3, 0x42, 0xb6, 0xea, 0xb0, 0xde, 0x56, 0xc9, 0x9a}, ValidationOpt: Enforce}, 6: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 8: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 9: {Expected: []byte{0xad, 0xca, 0x83, 0x3b, 0x2e, 0x4b, 0x6e, 0x91, 0x70, 0x89, 0xb9, 0xb5, 0xf3, 0x85, 0x08, 0x49, 0x91, 0x58, 0x6f, 0x92, 0xac, 0x2f, 0x69, 0x9c, 0x57, 0xf9, 0x0d, 0xeb, 0x88, 0xae, 0x7b, 0x85}, ValidationOpt: Enforce}, 11: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 12: {Expected: []byte{0xae, 0x8a, 0x39, 0x24, 0x54, 0x60, 0x84, 0x24, 0x17, 0x7b, 0x04, 0x1e, 0xc1, 0x5a, 0xc7, 0x25, 0x3d, 0x7e, 0x78, 0x0c, 0xed, 0x05, 0x90, 0xe1, 0x33, 0x2f, 0x0d, 0xdc, 0xae, 0xa2, 0xff, 0x4f}, ValidationOpt: Enforce}, 13: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 14: {Expected: []byte{0xd7, 0xc4, 0xcc, 0x7f, 0xf7, 0x93, 0x30, 0x22, 0xf0, 0x13, 0xe0, 0x3b, 0xde, 0xe8, 0x75, 0xb9, 0x17, 0x20, 0xb5, 0xb8, 0x6c, 0xf1, 0x75, 0x3c, 0xad, 0x83, 0x0f, 0x95, 0xe7, 0x91, 0x92, 0x6f}, ValidationOpt: WarnOnly}, 15: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}}
	azure_AzureSEVSNP        = M{1: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 2: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 3: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 4: {Expected: []byte{0x6e, 0x83, 0x6f, 0xd2, 0x8d, 0x6d, 0x74, 0xe7, 0x31, 0x1a, 0x60, 0xad, 0xaf, 0x44, 0x29, 0x71, 0xa8, 0xed, 0xba, 0xc0, 0x25, 0xc5, 0xc4, 0xfa, 0x69, 0xcf, 0xa4, 0xd2, 0x22, 0x7e, 0x06, 0x7f}, ValidationOpt: Enforce}, 8: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 9: {Expected: []byte{0xad, 0xca, 0x83, 0x3b, 0x2e, 0x4b, 0x6e, 0x91, 0x70, 0x89, 0xb9, 0xb5, 0xf3, 0x85, 0x08, 0x49, 0x91, 0x58, 0x6f, 0x92, 0xac, 0x2f, 0x69, 0x9c, 0x57, 0xf9, 0x0d, 0xeb, 0x88, 0xae, 0x7b, 0x85}, ValidationOpt: Enforce}, 11: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 12: {Expected: []byte{0xe4, 0x0d, 0x0c, 0xd1, 0x1f, 0x22, 0x33, 0x4d, 0xaa, 0xd4, 0x2c, 0x08, 0xd1, 0xfe, 0xda, 0x3b, 0xb6, 0x0d, 0x26, 0xba, 0x1a, 0x2a, 0xef, 0xa4, 0x7d, 0x9e, 0x94, 0xe6, 0x15, 0xce, 0xfa, 0x4c}, ValidationOpt: Enforce}, 13: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 14: {Expected: []byte{0xd7, 0xc4, 0xcc, 0x7f, 0xf7, 0x93, 0x30, 0x22, 0xf0, 0x13, 0xe0, 0x3b, 0xde, 0xe8, 0x75, 0xb9, 0x17, 0x20, 0xb5, 0xb8, 0x6c, 0xf1, 0x75, 0x3c, 0xad, 0x83, 0x0f, 0x95, 0xe7, 0x91, 0x92, 0x6f}, ValidationOpt: WarnOnly}, 15: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}}
	azure_AzureTrustedLaunch M
	gcp_GCPSEVES             = M{1: {Expected: []byte{0x74, 0x5f, 0x2f, 0xb4, 0x23, 0x5e, 0x46, 0x47, 0xaa, 0x0a, 0xd5, 0xac, 0xe7, 0x81, 0xcd, 0x92, 0x9e, 0xb6, 0x8c, 0x28, 0x87, 0x0e, 0x7d, 0xd5, 0xd1, 0xa1, 0x53, 0x58, 0x54, 0x32, 0x5e, 0x56}, ValidationOpt: WarnOnly}, 2: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 3: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 4: {Expected: []byte{0x8d, 0xb0, 0xe8, 0x59, 0x53, 0x0f, 0x39, 0x7c, 0x4d, 0x69, 0xfd, 0xa8, 0x4a, 0x4f, 0x29, 0xb4, 0x55, 0x05, 0x91, 0xd8, 0x69, 0x0b, 0xed, 0x09, 0xae, 0x9a, 0x00, 0x8e, 0x94, 0x95, 0x67, 0x77}, ValidationOpt: Enforce}, 6: {Expected: []byte{0x3d, 0x45, 0x8c, 0xfe, 0x55, 0xcc, 0x03, 0xea, 0x1f, 0x44, 0x3f, 0x15, 0x62, 0xbe, 0xec, 0x8d, 0xf5, 0x1c, 0x75, 0xe1, 0x4a, 0x9f, 0xcf, 0x9a, 0x72, 0x34, 0xa1, 0x3f, 0x19, 0x8e, 0x79, 0x69}, ValidationOpt: WarnOnly}, 8: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 9: {Expected: []byte{0xad, 0xca, 0x83, 0x3b, 0x2e, 0x4b, 0x6e, 0x91, 0x70, 0x89, 0xb9, 0xb5, 0xf3, 0x85, 0x08, 0x49, 0x91, 0x58, 0x6f, 0x92, 0xac, 0x2f, 0x69, 0x9c, 0x57, 0xf9, 0x0d, 0xeb, 0x88, 0xae, 0x7b, 0x85}, ValidationOpt: Enforce}, 11: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 12: {Expected: []byte{0xbd, 0x3a, 0xa0, 0x97, 0x44, 0xc8, 0xf7, 0xe8, 0xb4, 0x28, 0xb5, 0xa9, 0xaf, 0x10, 0x01, 0x3c, 0x06, 0x4d, 0xc5, 0xab, 0x50, 0xa4, 0x17, 0x60, 0xe6, 0x66, 0x37, 0x3d, 0x10, 0xfc, 0x62, 0xef}, ValidationOpt: Enforce}, 13: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 14: {Expected: []byte{0xd7, 0xc4, 0xcc, 0x7f, 0xf7, 0x93, 0x30, 0x22, 0xf0, 0x13, 0xe0, 0x3b, 0xde, 0xe8, 0x75, 0xb9, 0x17, 0x20, 0xb5, 0xb8, 0x6c, 0xf1, 0x75, 0x3c, 0xad, 0x83, 0x0f, 0x95, 0xe7, 0x91, 0x92, 0x6f}, ValidationOpt: WarnOnly}, 15: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}}
	qemu_QEMUTDX             M
	qemu_QEMUVTPM            = M{4: {Expected: []byte{0x59, 0x5a, 0xc1, 0x9e, 0x3b, 0xab, 0x6a, 0x13, 0x50, 0x48, 0x83, 0xaf, 0x3e, 0xb4, 0x0e, 0x65, 0x6f, 0xb1, 0x6c, 0xfe, 0x1a, 0x18, 0x84, 0x4e, 0x11, 0xe2, 0x6c, 0x0b, 0xd8, 0x2c, 0xad, 0xff}, ValidationOpt: Enforce}, 8: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 9: {Expected: []byte{0xad, 0xca, 0x83, 0x3b, 0x2e, 0x4b, 0x6e, 0x91, 0x70, 0x89, 0xb9, 0xb5, 0xf3, 0x85, 0x08, 0x49, 0x91, 0x58, 0x6f, 0x92, 0xac, 0x2f, 0x69, 0x9c, 0x57, 0xf9, 0x0d, 0xeb, 0x88, 0xae, 0x7b, 0x85}, ValidationOpt: Enforce}, 11: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 12: {Expected: []byte{0x10, 0x37, 0x7d, 0x82, 0x0d, 0x94, 0x34, 0x41, 0x6c, 0xa6, 0xec, 0x53, 0xa9, 0x7b, 0xc3, 0x12, 0xd5, 0x40, 0xb9, 0x06, 0xa2, 0x47, 0x0c, 0x6a, 0x56, 0xeb, 0x10, 0xf9, 0xfa, 0xec, 0x75, 0xce}, ValidationOpt: Enforce}, 13: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}, 15: {Expected: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, ValidationOpt: Enforce}}
)
