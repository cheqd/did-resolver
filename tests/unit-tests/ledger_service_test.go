package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

var _ = Describe("Test QueryDIDDoc method", func() {
	type testCase struct {
		did                        string
		expectedDidDocWithMetadata *didTypes.DidDocWithMetadata
		expectedIsFound            bool
		expectedError              error
	}

	It("cannot get DIDDoc with a invalid DID", func() {
		test := testCase{
			did:                        InvalidDid,
			expectedDidDocWithMetadata: nil,
			expectedIsFound:            false,
			expectedError:              types.NewInvalidDIDError(InvalidDid, types.JSON, nil, false),
		}

		ledgerService := services.NewLedgerService()
		didDocWithMetadata, err := ledgerService.QueryDIDDoc(InvalidDid, "")
		Expect(test.expectedDidDocWithMetadata).To(Equal(didDocWithMetadata))
		Expect(test.expectedError.Error()).To(Equal(err.Error()))
	})
})

var _ = Describe("Test QueryResource method", func() {
	type testCase struct {
		collectionId     string
		resourceId       string
		expectedResource *resourceTypes.ResourceWithMetadata
		expectedError    error
	}

	It("cannot get DIDDoc's resource with a invalid collectionId and resourceID", func() {
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
