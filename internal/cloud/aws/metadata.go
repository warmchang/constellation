/*
Copyright (c) Edgeless Systems GmbH

SPDX-License-Identifier: AGPL-3.0-only
*/

package aws

import (
	"context"

	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/edgelesssys/constellation/v2/internal/cloud/metadata"
)

type Metadata struct {
	ec2IAMInfo
	ec2InstanceIdentityDocument
	ec2Metadata
	ec2DynamicData
}

// NewMetadata creates a new Provider with the given API.
// Implement all functions defined in the ProviderMetadata interface.
func NewMetadata(ctx context.Context) (*Metadata, error) {
	ec2Data := ec2metadata.New(session.Must(session.NewSession()))

	IAMInfo, err := ec2Data.IAMInfoWithContext(ctx)
	if err != nil {
		return nil, err
	}

	InstanceIdentityDocument, err := ec2Data.GetInstanceIdentityDocumentWithContext(ctx)
	if err != nil {
		return nil, err
	}

	EC2Metadata, err := ec2Data.GetMetadataWithContext(ctx, "")
	if err != nil {
		return nil, err
	}

	DynamicData, err := ec2Data.GetDynamicDataWithContext(ctx, "")
	if err != nil {
		return nil, err
	}

	return &Metadata{
		ec2IAMInfo:                  &IAMInfo,
		ec2InstanceIdentityDocument: &InstanceIdentityDocument,
		ec2Metadata:                 &EC2Metadata,
		ec2DynamicData:              &DynamicData,
	}, nil
}

// UID retrieves the current instances uid.
func (m *Metadata) UID(ctx context.Context) (string, error) {
	return "", nil
}

// List retrieves all instances belonging to the current constellation.
func (m *Metadata) List(ctx context.Context) ([]metadata.InstanceMetadata, error) {
	return nil, nil
}

// Self retrieves the the current instance.
func (m *Metadata) Self(ctx context.Context) (metadata.InstanceMetadata, error) {
	return metadata.InstanceMetadata{}, nil
}

// GetSubnetworkCIDR retrieves the subnetwork CIDR of the current instance.
func (m *Metadata) GetSubnetworkCIDR(ctx context.Context) (string, error) {
	return "", nil
}

// SupportsLoadBalancers returns true if the cloud provider supports load balancers.
func (m *Metadata) SupportsLoadBalancer() bool {
	return false
}

// GetLoadBalancerEndpoint retrieves the load balancer endpoint of the current instance.
func (m *Metadata) GetLoadBalancerEndpoint(ctx context.Context) (string, error) {
	return "", nil
}

// GetInstance retrieves the instance type of the current instance.
func (m *Metadata) GetInstance(ctx context.Context, providerID string) (metadata.InstanceMetadata, error) {
	return metadata.InstanceMetadata{}, nil
}

// Supported is used to determine if metadata API is implemented for this cloud provider.
func (m *Metadata) Supported() bool {
	return false
}
