package tests

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase TestCase) {
	resourceService := services.NewResourceService(ValidMethod, mockLedgerService)
	id := "did:" + testCase.method + ":" + testCase.namespace + ":" + testCase.identifier

	var expectedDIDProperties types.DidProperties
	if testCase.expectedError == nil {
		expectedDIDProperties = types.DidProperties{
			DidString:        ValidDid,
			MethodSpecificId: ValidIdentifier,
			Method:           ValidMethod,
		}
	}

	expectedContentType := validResource.Metadata.MediaType
	dereferencingResult, err := resourceService.DereferenceResourceData(id, testCase.resourceId, testCase.dereferencingType)
	if err == nil {
		Expect(testCase.expectedContentStream, dereferencingResult.ContentStream)
		Expect(testCase.expectedMetadata, dereferencingResult.Metadata)
		Expect(expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
		Expect(expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
	} else {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	}
},

	Entry(
		"successful dereferencing for resource",
		TestCase{
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            ValidResourceId,
			expectedContentStream: &resourceData,
			expectedMetadata:      types.ResolutionResourceMetadata{},
			expectedError:         nil,
		},
	),

	Entry(
		"successful dereferencing for resource (upper case UUID)",
		TestCase{
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            strings.ToUpper(ValidResourceId),
			expectedContentStream: &resourceData,
			expectedMetadata:      types.ResolutionResourceMetadata{},
			expectedError:         nil,
		},
	),

	Entry(
		"resource not found",
		TestCase{
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        NotExistIdentifier,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid resource id",
		TestCase{
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        InvalidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"invalid method",
		TestCase{
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            InvalidMethod,
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier), types.DIDJSONLD, nil, true,
			),
		},
	),

	Entry(
		"invalid namespace",
		TestCase{
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         InvalidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier), types.DIDJSONLD, nil, true,
			),
		},
	),

	Entry(
		"invalid identifier",
		TestCase{
			dereferencingType: types.DIDJSONLD,
			identifier:        InvalidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionResourceMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier), types.DIDJSONLD, nil, true,
			),
		},
	),
)
