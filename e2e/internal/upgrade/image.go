//go:build e2e

/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package upgrade

import (
	"context"
	"errors"
	"net/http"

	"github.com/edgelesssys/constellation/v2/internal/attestation/measurements"
	"github.com/edgelesssys/constellation/v2/internal/cloud/cloudprovider"
	"github.com/edgelesssys/constellation/v2/internal/constants"
	"github.com/edgelesssys/constellation/v2/internal/versionsapi"
	"github.com/edgelesssys/constellation/v2/internal/versionsapi/fetcher"
)

type upgradeInfo struct {
	measurements measurements.M
	shortPath    string
	wantImage    string
}

func fetchUpgradeInfo(ctx context.Context, csp cloudprovider.Provider, toImage string) (upgradeInfo, error) {
	info := upgradeInfo{
		measurements: make(measurements.M),
		shortPath:    toImage,
	}
	versionsClient := fetcher.NewFetcher()

	ver, err := versionsapi.NewVersionFromShortPath(toImage, versionsapi.VersionKindImage)
	if err != nil {
		return upgradeInfo{}, err
	}

	measurementsURL, signatureURL, err := versionsapi.MeasurementURL(ver, csp)
	if err != nil {
		return upgradeInfo{}, err
	}

	var fetchedMeasurements measurements.M
	_, err = fetchedMeasurements.FetchAndVerify(
		ctx, http.DefaultClient,
		measurementsURL,
		signatureURL,
		[]byte(constants.CosignPublicKey),
		measurements.WithMetadata{
			CSP:   csp,
			Image: toImage,
		},
	)
	if err != nil {
		return upgradeInfo{}, err
	}
	info.measurements = fetchedMeasurements

	wantImage, err := fetchWantImage(ctx, versionsClient, csp, versionsapi.ImageInfo{
		Ref:     ver.Ref,
		Stream:  ver.Stream,
		Version: ver.Version,
	})
	if err != nil {
		return upgradeInfo{}, err
	}
	info.wantImage = wantImage

	return info, nil
}

func fetchWantImage(ctx context.Context, client *fetcher.Fetcher, csp cloudprovider.Provider, imageInfo versionsapi.ImageInfo) (string, error) {
	imageInfo, err := client.FetchImageInfo(ctx, imageInfo)
	if err != nil {
		return "", err
	}

	switch csp {
	case cloudprovider.GCP:
		return imageInfo.GCP["sev-es"], nil
	case cloudprovider.Azure:
		return imageInfo.Azure["cvm"], nil
	case cloudprovider.AWS:
		return imageInfo.AWS["eu-central-1"], nil
	default:
		return "", errors.New("finding wanted image")
	}
}