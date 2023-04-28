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
		"can get resource data with an existent collectionId and resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.ExistentDid,
			resourceId:       testconstants.ExistentResourceId,
			expectedResource: &testconstants.ValidResource[0],
			expectedError:    nil,
		},
	),

	Entry(
		"cannot get resource data with an existent collectionId, but not existent resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.ExistentDid,
			resourceId:       testconstants.NotExistentIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ExistentDid, types.JSON, nil, true),
		},
	),

	Entry(
		"cannot get resource data with not existent collectionId, but existent resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.NotExistentTestnetDid,
			resourceId:       testconstants.ExistentResourceId,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),

	Entry(
		"cannot get resource data with not existent collectionId and resourceId",
		queryResourceTestCase{
			collectionId:     testconstants.NotExistentTestnetDid,
			resourceId:       testconstants.NotExistentIdentifier,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),
)
