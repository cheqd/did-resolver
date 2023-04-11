package tests

import (
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type queryResourceTestCase struct {
	collectionId     string
	resourceId       string
	expectedResource *resourceTypes.ResourceWithMetadata
	expectedError    *types.IdentityError
}

var _ = DescribeTable("Test QueryDIDDoc method", func(testCase queryResourceTestCase) {
	resource, err := mockLedgerService.QueryResource(testCase.collectionId, testCase.resourceId)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(testCase.expectedResource).To(Equal(resource))
	}
},

	Entry(
		"existent collectionId and resourceId",
		queryResourceTestCase{
			collectionId:     ValidDid,
			resourceId:       ValidResourceId,
			expectedResource: &validResource,
			expectedError:    nil,
		},
	),

	Entry(
		"existent collectionId, but not existent resourceId",
		queryResourceTestCase{
			collectionId:     ValidDid,
			resourceId:       NotExistIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(ValidDid, types.JSON, nil, true),
		},
	),

	Entry(
		"existent collectionId, but an invalid resourceId",
		queryResourceTestCase{
			collectionId:     ValidDid,
			resourceId:       InvalidIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(ValidDid, types.JSON, nil, true),
		},
	),

	Entry(
		"not existent collectionId, but existent resourceId",
		queryResourceTestCase{
			collectionId:     NotExistDID,
			resourceId:       ValidResourceId,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(NotExistDID, types.JSON, nil, true),
		},
	),

	Entry(
		"an invalid collectionId, but existent resourceId",
		queryResourceTestCase{
			collectionId:     InvalidDid,
			resourceId:       ValidResourceId,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(InvalidDid, types.JSON, nil, true),
		},
	),

	Entry(
		"not existent collectionId and resourceId",
		queryResourceTestCase{
			collectionId:     NotExistDID,
			resourceId:       NotExistIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(NotExistDID, types.JSON, nil, true),
		},
	),
)
