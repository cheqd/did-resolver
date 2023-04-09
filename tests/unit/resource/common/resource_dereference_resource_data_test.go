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

type dereferenceResourceDataTestCase struct {
	did                           string
	resourceId                    string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
	resourceService := services.NewResourceService(testconstants.ValidMethod, utils.MockLedger)

	expectedContentType := types.ContentType(testconstants.ValidResource.Metadata.MediaType)
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
			did:               testconstants.ValidDid,
			resourceId:        testconstants.ValidResourceId,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &testconstants.ValidResourceDereferencing,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"successful dereferencing for resource (upper case UUID)",
		dereferenceResourceDataTestCase{
			did:               testconstants.ValidDid,
			resourceId:        strings.ToUpper(testconstants.ValidResourceId),
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &testconstants.ValidResourceDereferencing,
				Metadata:      types.ResolutionResourceMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"not existent DID and a valid resourceId",
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
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"an existent DID, but not existent resourceId",
		dereferenceResourceDataTestCase{
			did:               testconstants.ValidDid,
			resourceId:        testconstants.ValidIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"not existent DID and resourceId",
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
				Metadata: types.ResolutionResourceMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, true),
		},
	),
)
