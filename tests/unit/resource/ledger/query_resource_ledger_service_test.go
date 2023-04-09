//go:build unit

package ledger

import (
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
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

var _ = DescribeTable("Test QueryResource method", func(testCase queryResourceTestCase) {
	resource, err := utils.MockLedger.QueryResource(testCase.collectionId, testCase.resourceId)
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
			collectionId:     testconstants.ValidDid,
			resourceId:       testconstants.ValidResourceId,
			expectedResource: &testconstants.ValidResource,
			expectedError:    nil,
		},
	),

	Entry(
		"existent collectionId, but not existent resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.ValidDid,
			resourceId:       testconstants.NotExistentIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ValidDid, types.JSON, nil, true),
		},
	),

	Entry(
		"not existent collectionId, but existent resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.NotExistentTestnetDid,
			resourceId:       testconstants.ValidResourceId,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),

	Entry(
		"not existent collectionId and resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.NotExistentTestnetDid,
			resourceId:       testconstants.NotExistentIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),
)
