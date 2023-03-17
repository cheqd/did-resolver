package tests

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type TestCase struct {
	ledgerService         MockLedgerService
	dereferencingType     types.ContentType
	identifier            string
	method                string
	namespace             string
	resourceId            string
	expectedContentStream types.ContentStreamI
	expectedContentType   types.ContentType
	expectedMetadata      types.ResolutionResourceMetadata
	expectedError         *types.IdentityError
}

var _ = DescribeTable("Test DereferenceResourceMetadata method", func(testCase TestCase) {
	resourceService := services.NewResourceService(ValidMethod, testCase.ledgerService)
	id := "did:" + testCase.method + ":" + testCase.namespace + ":" + testCase.identifier

	var expectedDIDProperties types.DidProperties
	if testCase.expectedError == nil {
		expectedDIDProperties = types.DidProperties{
			DidString:        ValidDid,
			MethodSpecificId: ValidIdentifier,
			Method:           ValidMethod,
		}
	}

	expectedContentType := testCase.expectedContentType
	if expectedContentType == "" {
		expectedContentType = testCase.dereferencingType
	}

	dereferencingResult, err := resourceService.DereferenceResourceMetadata(testCase.resourceId, id, testCase.dereferencingType)
	if err == nil {
		Expect(testCase.expectedContentStream).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedMetadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(expectedDIDProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	} else {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	}
},

	Entry(
		"successful dereferencing for resource",
		TestCase{
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            ValidResourceId,
			expectedContentStream: dereferencedResourceList,
			expectedMetadata:      types.ResolutionResourceMetadata{},
			expectedError:         nil,
		},
	),

	Entry(
		"successful dereferencing for resource (upper case UUID)",
		TestCase{
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            strings.ToUpper(ValidResourceId),
			expectedContentStream: dereferencedResourceList,
			expectedMetadata:      types.ResolutionResourceMetadata{},
			expectedError:         nil,
		},
	),

	Entry(
		"resource not found",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid resource id",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        "invalid-resource-id",
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid resource id",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        "invalid-resource-id",
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid namespace",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         "invalid-namespace",
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid method",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            "invalid-method",
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid identifier",
		TestCase{
			ledgerService:     NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        "invalid-identifier",
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
