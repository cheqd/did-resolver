package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceDataTestCase struct {
	did                           string
	resourceId                    string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
	resourceService := services.NewResourceService(ValidMethod, mockLedgerService)

	expectedContentType := types.ContentType(validResource.Metadata.MediaType)
	dereferencingResult, err := resourceService.DereferenceResourceData(testCase.did, testCase.resourceId, testCase.dereferencingType)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedResourceDereferencing.ContentStream.GetBytes()).To(Equal(dereferencingResult.ContentStream.GetBytes()))
		Expect(testCase.expectedResourceDereferencing.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedResourceDereferencing.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	}
},

	Entry(
		"successful dereferencing for resource",
		dereferenceResourceDataTestCase{
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
				ContentStream: &validResourceDereferencing,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"successful dereferencing for resource (upper case UUID)",
		dereferenceResourceDataTestCase{
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
				ContentStream: &validResourceDereferencing,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"resource not found",
		dereferenceResourceDataTestCase{
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
		"invalid resource id",
		dereferenceResourceDataTestCase{
			did:               ValidDid,
			resourceId:        InvalidResourceId,
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
		"invalid method",
		dereferenceResourceDataTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier),
			resourceId:        InvalidResourceId,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier),
				types.DIDJSONLD, nil, true,
			),
		},
	),

	Entry(
		"invalid namespace",
		dereferenceResourceDataTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier),
			resourceId:        InvalidResourceId,
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
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier),
				types.DIDJSONLD, nil, true,
			),
		},
	),

	Entry(
		"invalid identifier",
		dereferenceResourceDataTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier),
			resourceId:        InvalidResourceId,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        InvalidDid,
						MethodSpecificId: InvalidIdentifier,
						Method:           ValidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier),
				types.DIDJSONLD, nil, true,
			),
		},
	),

	// TODO: add unit tests for:
	// - invalid DID and invalid resourceId
	// - not existent DID and existent resourceId
	// - not existent DID and not existent resourceId.
)
