//go:build unit

package ledger

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type queryDIDDocTestCase struct {
	did                        string
	expectedDidDocWithMetadata *didTypes.DidDocWithMetadata
	expectedError              *types.IdentityError
}

var _ = DescribeTable("Test QueryDIDDoc method", func(testCase queryDIDDocTestCase) {
	didDocWithMetadata, err := utils.MockLedger.QueryDIDDoc(testCase.did, "")
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(testCase.expectedDidDocWithMetadata).To(Equal(didDocWithMetadata))
	}
},

	Entry(
		"can get DIDDoc with an existent DID",
		queryDIDDocTestCase{
			did: testconstants.ExistentDid,
			expectedDidDocWithMetadata: &didTypes.DidDocWithMetadata{
				DidDoc:   &testconstants.ValidDIDDoc,
				Metadata: &testconstants.ValidMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get DIDDoc with not existent DID",
		queryDIDDocTestCase{
			did:                        testconstants.NotExistentTestnetDid,
			expectedDidDocWithMetadata: nil,
			expectedError:              types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),
)
