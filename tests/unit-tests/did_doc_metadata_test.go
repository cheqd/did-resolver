package tests

import (
	"crypto/sha256"
	"fmt"
	"testing"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
	// "google.golang.org/protobuf/types/known/timestamppb"
)

func TestNewResolutionDidDocMetadata(t *testing.T) {
	validIdentifier := "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
	validDid := "did:cheqd:mainnet:" + validIdentifier
	validResourceId := "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
	resourceData := []byte("test_checksum")
	h := sha256.New()
	h.Write(resourceData)
	resourceMetadata := resourceTypes.Metadata{
		CollectionId: validIdentifier,
		Id:           validResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     fmt.Sprintf("%x", h.Sum(nil)),
	}

	created := resourceMetadata.Created.AsTime()
	validMetadataResource := types.DereferencedResource{
		ResourceURI:       validDid + types.RESOURCE_PATH + resourceMetadata.Id,
		CollectionId:      resourceMetadata.CollectionId,
		ResourceId:        resourceMetadata.Id,
		Name:              resourceMetadata.Name,
		ResourceType:      resourceMetadata.ResourceType,
		MediaType:         resourceMetadata.MediaType,
		Created:           &created,
		Checksum:          resourceMetadata.Checksum,
		PreviousVersionId: nil,
		NextVersionId:     nil,
	}

	subtests := []struct {
		name           string
		metadata       *didTypes.Metadata
		resources      []*resourceTypes.Metadata
		expectedResult types.ResolutionDidDocMetadata
	}{
		{
			name: "metadata with resource",
			metadata: &didTypes.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			resources: []*resourceTypes.Metadata{&resourceMetadata},
			expectedResult: types.ResolutionDidDocMetadata{
				Created:     &EmptyTime,
				Updated:     &EmptyTime,
				Deactivated: false,
				VersionId:   "test_version_id",
				Resources:   []types.DereferencedResource{validMetadataResource},
			},
		},
		{
			name: "metadata without resources",
			metadata: &didTypes.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			expectedResult: types.ResolutionDidDocMetadata{
				Created:     &EmptyTime,
				Updated:     &EmptyTime,
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
