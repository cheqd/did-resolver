//go:build unit

package ledger

import (
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type queryDIDDocVersionsTestCase struct {
	did                    string
	expectedDidDocVersions *types.DereferencedDidVersionsList
	expectedError          *types.IdentityError
}

var _ = DescribeTable("Test QueryAllDidDocVersionsMetadata method", func(testCase queryDIDDocVersionsTestCase) {
	didDocMetadata, err := utils.MockLedger.QueryAllDidDocVersionsMetadata(testCase.did)
	didDocVersions := types.NewDereferencedDidVersionsList(didDocMetadata)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(testCase.expectedDidDocVersions).To(Equal(didDocVersions))
	}
},

	Entry(
		"can get DIDDoc versions with an existent DID",
		queryDIDDocVersionsTestCase{
			did: testconstants.ExistentDid,
			expectedDidDocVersions: &types.DereferencedDidVersionsList{
				Versions: []types.ResolutionDidDocMetadata{
					{
						VersionId: testconstants.ValidMetadata.VersionId,
					},
				},
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get DIDDoc versions with not existent DID",
		queryDIDDocVersionsTestCase{
			did:                    testconstants.NotExistentTestnetDid,
			expectedDidDocVersions: nil,
			expectedError:          types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),
)
