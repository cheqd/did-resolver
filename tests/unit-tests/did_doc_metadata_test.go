package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
)

var _ = Describe("Test NewResolutionDIDDocMetadata function", func() {
	It("can resolution metadata with resource", func() {
		metadata := &didTypes.Metadata{
			VersionId:   "test_version_id",
			Deactivated: false,
		}

		resources := []*resourceTypes.Metadata{
			&ResourceMetadata,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:     nil,
			Updated:     nil,
			Deactivated: false,
			VersionId:   "test_version_id",
			Resources:   []types.DereferencedResource{ValidMetadataResource},
		}

		result := types.NewResolutionDidDocMetadata(ValidDid, metadata, resources)
		Expect(result).To(Equal(expectedResult))
	})

	It("can resolution metadata without resource", func() {
		metadata := &didTypes.Metadata{
			Created:     NotEmptyTimestamp,
			Updated:     NotEmptyTimestamp,
			VersionId:   "test_version_id",
			Deactivated: false,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:     &NotEmptyTime,
			Updated:     &NotEmptyTime,
			VersionId:   "test_version_id",
			Deactivated: false,
		}

		result := types.NewResolutionDidDocMetadata(ValidDid, metadata, []*resourceTypes.Metadata{})
		Expect(result).To(Equal(expectedResult))
	})
})
