package types

import (
	"crypto/sha256"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/stretchr/testify/require"
)

func TestNewResolutionDidDocMetadata(t *testing.T) {
	validIdentifier := "N22KY2Dyvmuu2Pyy"
	validDid := "did:cheqd:mainnet:" + validIdentifier
	validResourceId := "18e9d838-0bea-435b-964b-c6529ede6d2b"
	resourceData := []byte("test_checksum")
	h := sha256.New()
	h.Write(resourceData)
	resourceHeader := resource.ResourceHeader{
		CollectionId: validIdentifier,
		Id:           validResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     h.Sum(nil),
	}

	validMetadataResource := DereferencedResource{
		ResourceURI:       validDid + RESOURCE_PATH + resourceHeader.Id,
		CollectionId:      resourceHeader.CollectionId,
		ResourceId:        resourceHeader.Id,
		Name:              resourceHeader.Name,
		ResourceType:      resourceHeader.ResourceType,
		MediaType:         resourceHeader.MediaType,
		Created:           resourceHeader.Created,
		Checksum:          FixResourceChecksum(resourceHeader.Checksum),
		PreviousVersionId: nil,
		NextVersionId:     nil,
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
				Resources:   []DereferencedResource{validMetadataResource},
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
			name: "matadata without resources",
			metadata: cheqd.Metadata{
				VersionId:   "test_version_id",
				Deactivated: false,
				Resources:   []string{validResourceId},
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
