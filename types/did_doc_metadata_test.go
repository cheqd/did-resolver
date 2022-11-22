package types

import (
	"crypto/sha256"
	"fmt"
	"testing"

	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/stretchr/testify/require"
)

func TestNewResolutionDidDocMetadata(t *testing.T) {
	validIdentifier := "0a88fd42-ae6b-4da4-a5ea-36d28e35dde7"
	validDid := "did:cheqd:mainnet:" + validIdentifier
	validResourceId := "18e9d838-0bea-435b-964b-c6529ede6d2b"
	resourceData := []byte("test_checksum")
	h := sha256.New()
	h.Write(resourceData)
	resourceMetadata := resource.Metadata{
		CollectionId: validIdentifier,
		Id:           validResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     h.Sum(nil),
	}

	validMetadataResource := DereferencedResource{
		ResourceURI:       validDid + RESOURCE_PATH + resourceMetadata.Id,
		CollectionId:      resourceMetadata.CollectionId,
		ResourceId:        resourceMetadata.Id,
		Name:              resourceMetadata.Name,
		ResourceType:      resourceMetadata.ResourceType,
		MediaType:         resourceMetadata.MediaType,
		Created:           resourceMetadata.Created,
		Checksum:          fmt.Sprintf("%x", resourceMetadata.Checksum),
		PreviousVersionId: nil,
		NextVersionId:     nil,
	}

	subtests := []struct {
		name           string
		metadata       did.Metadata
		resources      []*resource.Metadata
		expectedResult ResolutionDidDocMetadata
	}{
		{
			name: "matadata with resource",
			metadata: did.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			resources: []*resource.Metadata{&resourceMetadata},
			expectedResult: ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []DereferencedResource{validMetadataResource},
			},
		},
		{
			name: "matadata without resource in metadata",
			metadata: did.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			resources: []*resource.Metadata{&resourceMetadata},
			expectedResult: ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
		},
		{
			name: "matadata without resources",
			metadata: did.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			expectedResult: ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			result := NewResolutionDidDocMetadata(validDid, subtest.metadata, subtest.resources)

			require.EqualValues(t, subtest.expectedResult, result)
		})
	}
}
