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

type queryCollectionResourcesTestCase struct {
	did                string
	expectedCollection []*resourceTypes.Metadata
	expectedError      *types.IdentityError
}

var _ = DescribeTable("Test something", func(testCase queryCollectionResourcesTestCase) {
	collection, err := utils.MockLedger.QueryCollectionResources(testCase.did)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(testCase.expectedCollection).To(Equal(collection))
	}
},

	Entry(
		"can get collection of resources an existent DID",
		queryCollectionResourcesTestCase{
			did: testconstants.ExistentDid,
			expectedCollection: []*resourceTypes.Metadata{
				testconstants.ValidResource.Metadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"can get collection of resources not existent DID",
		queryCollectionResourcesTestCase{
			did:                testconstants.NotExistentTestnetDid,
			expectedCollection: nil,
			expectedError:      types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.JSON, nil, true),
		},
	),
)
