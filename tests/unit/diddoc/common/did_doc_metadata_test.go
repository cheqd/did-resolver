//go:build unit

package common

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
)

var _ = Describe("Test NewResolutionDIDDocMetadata function", func() {
	It("can create the structure with resource", func() {
		metadata := &didTypes.Metadata{
			VersionId:   testconstants.ValidIdentifier,
			Deactivated: false,
		}

		resources := []*resourceTypes.Metadata{
			&testconstants.ValidResourceMetadata,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:     nil,
			Updated:     nil,
			Deactivated: false,
			VersionId:   testconstants.ValidIdentifier,
			Resources:   []types.DereferencedResource{testconstants.ValidMetadataResource},
		}

		result := types.NewResolutionDidDocMetadata(testconstants.ValidDid, metadata, resources)
		Expect(result).To(Equal(expectedResult))
	})

	It("can create the structure without resource", func() {
		metadata := &didTypes.Metadata{
			Created:     testconstants.NotEmptyTimestamp,
			Updated:     testconstants.NotEmptyTimestamp,
			VersionId:   testconstants.ValidVersionId,
			Deactivated: false,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:     &testconstants.NotEmptyTime,
			Updated:     &testconstants.NotEmptyTime,
			VersionId:   testconstants.ValidVersionId,
			Deactivated: false,
		}

		result := types.NewResolutionDidDocMetadata(testconstants.ValidDid, metadata, []*resourceTypes.Metadata{})
		Expect(result).To(Equal(expectedResult))
	})
})
