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

type dereferenceResourceMetadataTestCase struct {
	did                           string
	resourceId                    string
	dereferencingType             types.ContentType
	expectedResourceDereferencing *types.ResourceDereferencing
	expectedError                 *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceMetadata method", func(testCase dereferenceResourceMetadataTestCase) {
	resourceService := services.NewResourceService(testconstants.ValidMethod, utils.MockLedger)

	expectedContentType := utils.DefineContentType(
		testCase.expectedResourceDereferencing.DereferencingMetadata.ContentType,
		testCase.dereferencingType,
	)

	dereferencingResult, err := resourceService.DereferenceResourceMetadata(testCase.did, testCase.resourceId, testCase.dereferencingType)
	if err != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(testCase.expectedResourceDereferencing.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedResourceDereferencing.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	}
},

	Entry(
		"can get resource metadata with an existent DID and resourceId",
		dereferenceResourceMetadataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        testconstants.ExistentResourceId,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Metadata: types.NewDereferencedResource(
					testconstants.ExistentDid,
					testconstants.ValidResource[0].Metadata,
				),
			},
			expectedError: nil,
		},
	),

	Entry(
		"can get resource metadata with an existent DID and upper case resourceId",
		dereferenceResourceMetadataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        strings.ToUpper(testconstants.ExistentResourceId),
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Metadata: types.NewDereferencedResource(
					testconstants.ExistentDid,
					testconstants.ValidResource[0].Metadata,
				),
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get resource metadata with not existent DID",
		dereferenceResourceMetadataTestCase{
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
				Metadata: types.NewDereferencedResource(
					testconstants.ExistentDid,
					testconstants.ValidResource[0].Metadata,
				),
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"cannot get resource metadata with an existent DID, but not existent resourceId",
		dereferenceResourceMetadataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        testconstants.NotExistentIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      nil,
			},
			expectedError: types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"cannot get resource metadata not existent DID and resourceId",
		dereferenceResourceMetadataTestCase{
			did:               testconstants.ExistentDid,
			resourceId:        testconstants.NotExistentIdentifier,
			dereferencingType: types.DIDJSON,
			expectedResourceDereferencing: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      nil,
			},
			expectedError: types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, true),
		},
	),
)
