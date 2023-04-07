package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	"github.com/cheqd/did-resolver/types"
)

type queryDIDDocTestCase struct {
	did                        string
	expectedDidDocWithMetadata *didTypes.DidDocWithMetadata
	expectedError              *types.IdentityError
}

var _ = DescribeTable("Test QueryDIDDoc method", func(testCase queryDIDDocTestCase) {
	didDocWithMetadata, err := mockLedgerService.QueryDIDDoc(testCase.did, "")
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(testCase.expectedDidDocWithMetadata).To(Equal(didDocWithMetadata))
	}
},

	Entry(
		"existent DID",
		queryDIDDocTestCase{
			did: ValidDid,
			expectedDidDocWithMetadata: &didTypes.DidDocWithMetadata{
				DidDoc:   &validDIDDoc,
				Metadata: &validMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"not existent DID",
		queryDIDDocTestCase{
			did:                        NotExistDID,
			expectedDidDocWithMetadata: nil,
			expectedError:              types.NewNotFoundError(NotExistDID, types.JSON, nil, true),
		},
	),
)

// TODO: add unit tests for testing other Ledger services:
// - QueryAllDidDocVersionsMetadata
// - QueryCollectionResources
