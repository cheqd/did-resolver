package types

import (
	"fmt"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/stretchr/testify/require"
)

func TestNewResolutionDidDocMetadata(t *testing.T) {
	validIdentifier := "N22KY2Dyvmuu2Pyy"
	validDid := "did:cheqd:mainnet:" + validIdentifier
	validResourceId := "18e9d838-0bea-435b-964b-c6529ede6d2b"
	resourceHeader := resource.ResourceHeader{
		CollectionId: validIdentifier,
		Id:           validResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     []byte("test_checksum"),
	}

	validMetadataResource := ResourcePreview{
		ResourceURI:       validDid + RESOURCE_PATH + resourceHeader.Id,
		Name:              resourceHeader.Name,
		ResourceType:      resourceHeader.ResourceType,
		MediaType:         resourceHeader.MediaType,
		Created:           resourceHeader.Created,
		Checksum:          fmt.Sprintf("%x", resourceHeader.Checksum),
		PreviousVersionId: resourceHeader.PreviousVersionId,
		NextVersionId:     resourceHeader.NextVersionId,
	}

	subtests := []struct {
		name           string
		metadata       cheqd.Metadata
		resources      []*resource.ResourceHeader
		expectedResult ResolutionDidDocMetadata
	}{
		{
			name: "matadata with resource",
			metadata: cheqd.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []string{validResourceId},
			},
			resources: []*resource.ResourceHeader{&resourceHeader},
			expectedResult: ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []ResourcePreview{validMetadataResource},
			},
		},
		{
			name: "matadata without resource in metadata",
			metadata: cheqd.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
			resources: []*resource.ResourceHeader{&resourceHeader},
			expectedResult: ResolutionDidDocMetadata{
				VersionId:   "test_version_id",
				Deactivated: false,
			},
		},
		{
			name: "matadata with resources",
			metadata: cheqd.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []string{validResourceId},
			},
			resources: []*resource.ResourceHeader{},
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
			// require.EqualValues(t, subtest.expectedError, err)
		})
	}
}
