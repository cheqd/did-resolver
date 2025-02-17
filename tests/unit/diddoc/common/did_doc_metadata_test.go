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
			Created:     nil,
		}

		resources := []*resourceTypes.Metadata{
			&testconstants.ValidResourceMetadata,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:           nil,
			Updated:           nil,
			Deactivated:       false,
			NextVersionId:     "",
			PreviousVersionId: "",
			VersionId:         testconstants.ValidIdentifier,
			Resources:         []types.DereferencedResource{testconstants.ValidMetadataResource},
		}

		result := types.NewResolutionDidDocMetadata(testconstants.ExistentDid, metadata, resources)
		Expect(expectedResult).To(Equal(result))
	})

	It("can create the structure without resource", func() {
		metadata := &didTypes.Metadata{
			Created:           testconstants.NotEmptyTimestamp,
			Updated:           testconstants.NotEmptyTimestamp,
			VersionId:         testconstants.ValidVersionId,
			NextVersionId:     testconstants.ValidNextVersionId,
			PreviousVersionId: testconstants.ValidPreviousVersionId,
			Deactivated:       false,
		}

		expectedResult := types.ResolutionDidDocMetadata{
			Created:           &testconstants.NotEmptyTime,
			Updated:           &testconstants.NotEmptyTime,
			VersionId:         testconstants.ValidVersionId,
			NextVersionId:     testconstants.ValidNextVersionId,
			PreviousVersionId: testconstants.ValidPreviousVersionId,
			Deactivated:       false,
		}

		result := types.NewResolutionDidDocMetadata(testconstants.ExistentDid, metadata, []*resourceTypes.Metadata{})
		Expect(expectedResult).To(Equal(result))
	})
})
