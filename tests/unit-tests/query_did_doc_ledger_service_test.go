package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
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

	Entry(
		"invalid DID",
		queryDIDDocTestCase{
			did:                        InvalidDid,
			expectedDidDocWithMetadata: nil,
			expectedError:              types.NewNotFoundError(InvalidDid, types.JSON, nil, true),
		},
	),
)

var _ = Describe("Test QueryResource method", func() {
	type testCase struct {
		collectionId     string
		resourceId       string
		expectedResource *resourceTypes.ResourceWithMetadata
		expectedError    error
	}

	It("cannot get DIDDoc's resource with a invalid collectionId and resourceId", func() {
		test := testCase{
			collectionId:     InvalidDid,
			resourceId:       InvalidResourceId,
			expectedResource: nil,
			expectedError:    types.NewInvalidDIDError(InvalidDid, types.JSON, nil, true),
		}

		ledgerService := services.NewLedgerService()
		resource, err := ledgerService.QueryResource(test.collectionId, test.resourceId)
		Expect(test.expectedResource).To(Equal(resource))
		Expect(test.expectedError.Error()).To(Equal(err.Error()))
	})
})
