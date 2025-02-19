//go:build unit

package common

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceCollectionResourcesTestCase struct {
	did                           string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceCollectionResources method", func(testCase dereferenceCollectionResourcesTestCase) {
	resourceService := services.NewResourceService(testconstants.ValidMethod, utils.MockLedger)

	expectedContentType := utils.DefineContentType(
		testCase.expectedResourceDereferencing.DereferencingMetadata.ContentType,
		testCase.dereferencingType,
	)

	dereferencingResult, err := resourceService.DereferenceCollectionResources(testCase.did, testCase.dereferencingType)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedResourceDereferencing.ContentStream, dereferencingResult.ContentStream)
		Expect(testCase.expectedResourceDereferencing.Metadata, dereferencingResult.Metadata)
		Expect(expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
		Expect(testCase.expectedResourceDereferencing.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	}
},

	Entry(
		"can get collection of resources with an existent DID",
		dereferenceCollectionResourcesTestCase{
			did:               testconstants.ExistentDid,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: testconstants.ValidDereferencedResourceList,
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get collection of resources with not existent DID",
		dereferenceCollectionResourcesTestCase{
			did:               testconstants.NotExistentTestnetDid,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, true),
		},
	),
)
