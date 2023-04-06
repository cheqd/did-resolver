package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceCollectionResourcesTestCase struct {
	did                           string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceCollectionResources method", func(testCase dereferenceCollectionResourcesTestCase) {
	resourceService := services.NewResourceService(ValidMethod, mockLedgerService)

	expectedContentType := defineContentType(
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
		"successful dereferencing for collection resources",
		dereferenceCollectionResourcesTestCase{
			did:               ValidDid,
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
		"invalid method",
		dereferenceCollectionResourcesTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier),
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
		dereferenceCollectionResourcesTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier),
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
		dereferenceCollectionResourcesTestCase{
			did:               fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier),
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

	Entry(
		"invalid did",
		dereferenceCollectionResourcesTestCase{
			did:               InvalidDid,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        InvalidDid,
						MethodSpecificId: InvalidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(InvalidDid, types.DIDJSONLD, nil, true),
		},
	),
)
