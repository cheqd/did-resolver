//go:build unit

package common

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceDataWithMetadataTestCase struct {
	did                           string
	resourceId                    string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
	resourceService := services.NewResourceService(testconstants.ValidMethod, utils.MockLedger)

	expectedContentType := types.ContentType(testconstants.ValidResource[0].Metadata.MediaType)
	dereferencingResult, err := resourceService.DereferenceResourceDataWithMetadata(testCase.did, testCase.resourceId, testCase.dereferencingType)
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
		"can get resource data with an existent DID and resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        testconstants.ExistentResourceId,
			dereferencingType: types.JSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &testconstants.ValidResourceDereferencing,
				Metadata: &types.ResolutionResourceMetadata{
					ContentMetadata: types.NewDereferencedResource(
						testconstants.ExistentDid,
						testconstants.ValidResource[0].Metadata,
					),
				},
			},
			expectedError: nil,
		},
	),

	Entry(
		"can get resource data with an existent DID and upper case resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        strings.ToUpper(testconstants.ExistentResourceId),
			dereferencingType: types.JSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &testconstants.ValidResourceDereferencing,
				Metadata: &types.ResolutionResourceMetadata{
					ContentMetadata: types.NewDereferencedResource(
						testconstants.ExistentDid,
						testconstants.ValidResource[0].Metadata,
					),
				},
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and a valid resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.NotExistentTestnetDid,
			resourceId:        testconstants.ValidIdentifier,
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

	Entry(
		"cannot get resource data with an existent DID, but not existent resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        testconstants.ValidIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
			},
			expectedError: types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"cannot get resource data with not existent DID and resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.NotExistentTestnetDid,
			resourceId:        testconstants.NotExistentIdentifier,
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
