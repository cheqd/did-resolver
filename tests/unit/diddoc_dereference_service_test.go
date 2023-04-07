package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type dereferencingTestCase struct {
	did                      string
	fragmentId               string
	dereferencingType        types.ContentType
	expectedDidDereferencing *types.DidDereferencing
	expectedError            *types.IdentityError
}

var _ = DescribeTable("Test Dereferencing method", func(testCase dereferencingTestCase) {
	diddocService := services.NewDIDDocService("cheqd", mockLedgerService)

	expectedContentType := defineContentType(
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
		"successful Secondary dereferencing (verification method)",
		dereferencingTestCase{
			did:               ValidDid,
			fragmentId:        validVerificationMethod.Id,
			dereferencingType: types.DIDJSON,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: types.NewVerificationMethod(&validVerificationMethod),
				Metadata:      validFragmentMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"successful Secondary dereferencing (service)",
		dereferencingTestCase{
			did:               ValidDid,
			fragmentId:        validService.Id,
			dereferencingType: types.DIDJSON,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: types.NewService(&validService),
				Metadata:      validFragmentMetadata,
			},
			expectedError: nil,
		},
	),

	Entry(
		"key not found",
		dereferencingTestCase{
			did:               ValidDid,
			fragmentId:        NotExistFragmentId,
			dereferencingType: types.DIDJSONLD,
			expectedDidDereferencing: &types.DidDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)
