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

type dereferencingTestCase struct {
	did                      string
	fragmentId               string
	dereferencingType        types.ContentType
	expectedDidDereferencing *types.DidDereferencing
	expectedError            *types.IdentityError
}

var _ = DescribeTable("Test DereferenceSecondary method", func(testCase dereferencingTestCase) {
	diddocService := services.NewDIDDocService("cheqd", utils.MockLedger)

	expectedContentType := utils.DefineContentType(
		testCase.expectedDidDereferencing.DereferencingMetadata.ContentType, testCase.dereferencingType,
	)

	dereferencingResult, err := diddocService.DereferenceSecondary(testCase.did, "", testCase.fragmentId, testCase.dereferencingType)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedDidDereferencing.ContentStream).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedDidDereferencing.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
		Expect(testCase.expectedDidDereferencing.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
	}
},

	Entry(
		"can successful dereferencing secondary (verification method) with an existent DID and verificationMethodId",
		dereferencingTestCase{
			did:               testconstants.ExistentDid,
			fragmentId:        testconstants.ValidVerificationMethod.Id,
			dereferencingType: types.DIDJSON,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: types.NewVerificationMethod(&testconstants.ValidVerificationMethod),
				Metadata:      *testconstants.ValidFragmentMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"can successful dereferencing secondary (service) with an existent DID and serviceId",
		dereferencingTestCase{
			did:               testconstants.ExistentDid,
			fragmentId:        testconstants.ValidService.Id,
			dereferencingType: types.DIDJSON,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: types.NewService(&testconstants.ValidService),
				Metadata:      *testconstants.ValidFragmentMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot dereferencing secondary with an existent DID, but not existent fragment",
		dereferencingTestCase{
			did:               testconstants.ExistentDid,
			fragmentId:        testconstants.NotExistentFragment,
			dereferencingType: types.DIDJSONLD,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, false),
		},
	),
)
