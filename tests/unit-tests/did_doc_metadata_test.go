package tests

import (
	"crypto/sha256"
	"fmt"
	"testing"

	did "github.com/cheqd/cheqd-node/api/cheqd/did/v2"
	resource "github.com/cheqd/cheqd-node/api/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
)

func TestNewResolutionDidDocMetadata(t *testing.T) {
	validIdentifier := "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
	validDid := "did:cheqd:mainnet:" + validIdentifier
	validResourceId := "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
	resourceData := []byte("test_checksum")
	h := sha256.New()
	h.Write(resourceData)
	resourceMetadata := resource.Metadata{
		CollectionId: validIdentifier,
		Id:           validResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     fmt.Sprintf("%x", h.Sum(nil)),
	}

	validMetadataResource := types.DereferencedResource{
		ResourceURI:       validDid + types.RESOURCE_PATH + resourceMetadata.Id,
		CollectionId:      resourceMetadata.CollectionId,
		ResourceId:        resourceMetadata.Id,
		Name:              resourceMetadata.Name,
		ResourceType:      resourceMetadata.ResourceType,
		MediaType:         resourceMetadata.MediaType,
		Created:           resourceMetadata.Created,
		Checksum:          resourceMetadata.Checksum,
		PreviousVersionId: nil,
		NextVersionId:     nil,
	}

	subtests := []struct {
		name           string
		metadata       did.Metadata
		resources      []*resource.Metadata
		expectedResult types.ResolutionDidDocMetadata
	}{
		{
			name: "metadata with resource",
			metadata: did.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			resources: []*resource.Metadata{&resourceMetadata},
			expectedResult: types.ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []types.DereferencedResource{validMetadataResource},
			},
		},
		{
			name: "metadata without resources",
			metadata: did.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			expectedResult: types.ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			result := types.NewResolutionDidDocMetadata(validDid, subtest.metadata, subtest.resources)

			require.EqualValues(t, subtest.expectedResult, result)
		})
	}
}
