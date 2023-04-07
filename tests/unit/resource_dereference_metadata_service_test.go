package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceMetadataTestCase struct {
	did                           string
	resourceId                    string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceMetadata method", func(testCase dereferenceResourceMetadataTestCase) {
	resourceService := services.NewResourceService(ValidMethod, mockLedgerService)

	expectedContentType := defineContentType(
		testCase.expectedResourceDereferencing.DereferencingMetadata.ContentType,
		testCase.dereferencingType,
	)

	dereferencingResult, err := resourceService.DereferenceResourceMetadata(testCase.did, testCase.resourceId, testCase.dereferencingType)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedResourceDereferencing.ContentStream).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedResourceDereferencing.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedResourceDereferencing.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	}
},

	Entry(
		"successful dereferencing for resource",
		dereferenceResourceMetadataTestCase{
			did:               ValidDid,
			resourceId:        ValidResourceId,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: dereferencedResourceList,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"successful dereferencing for resource (upper case UUID)",
		dereferenceResourceMetadataTestCase{
			did:               ValidDid,
			resourceId:        strings.ToUpper(ValidResourceId),
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: dereferencedResourceList,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"not existent DID",
		dereferenceResourceMetadataTestCase{
			did:               NotExistDID,
			resourceId:        NotExistIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        NotExistDID,
						MethodSpecificId: NotExistIdentifier,
						Method:           ValidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"an existent DID, but not existent resourceId",
		dereferenceResourceMetadataTestCase{
			did:               ValidDid,
			resourceId:        NotExistIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"not existent DID and resourceId",
		dereferenceResourceMetadataTestCase{
			did:               ValidDid,
			resourceId:        NotExistIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
