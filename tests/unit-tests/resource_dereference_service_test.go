package tests

import (
	"fmt"
	"strings"
	"testing"

	did "github.com/cheqd/cheqd-node/api/cheqd/did/v2"
	resource "github.com/cheqd/cheqd-node/api/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name                  string
	ledgerService         MockLedgerService
	dereferencingType     types.ContentType
	identifier            string
	method                string
	namespace             string
	resourceId            string
	expectedContentStream types.ContentStreamI
	expectedContentType   types.ContentType
	expectedMetadata      types.ResolutionDidDocMetadata
	expectedError         *types.IdentityError
}

func getSubtest(validContentStream types.ContentStreamI) []TestCase {
	validDIDDoc := ValidDIDDoc()
	validResource := ValidResource()
	validMetadata := ValidMetadata()
	return []TestCase{
		{
			name:                  "successful dereferencing for resource",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            ValidResourceId,
			expectedContentStream: validContentStream,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         nil,
		},
		{
			name:                  "successful dereferencing for resource (upper case UUID)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            ValidIdentifier,
			method:                ValidMethod,
			namespace:             ValidNamespace,
			resourceId:            strings.ToUpper(ValidResourceId),
			expectedContentStream: validContentStream,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         nil,
		},
		{
			name:              "resource not found",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid resource id",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        "invalid-resource-id",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid type",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.JSON,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid namespace",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            ValidMethod,
			namespace:         "invalid-namespace",
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid method",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        ValidIdentifier,
			method:            "invalid-method",
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid identifier",
			ledgerService:     NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        "invalid-identifier",
			method:            ValidMethod,
			namespace:         ValidNamespace,
			resourceId:        ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(ValidDid, types.DIDJSONLD, nil, true),
		},
	}
}

func TestDereferenceResourceMetadata(t *testing.T) {
	validResource := ValidResource()
	subtests := getSubtest(types.NewDereferencedResourceList(ValidDid, []*resource.Metadata{validResource.Metadata}))

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			resourceService := services.NewResourceService(ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        ValidDid,
					MethodSpecificId: ValidIdentifier,
					Method:           ValidMethod,
				}
			}
			expectedContentType := subtest.expectedContentType
			if expectedContentType == "" {
				expectedContentType = subtest.dereferencingType
			}
			dereferencingResult, err := resourceService.DereferenceResourceMetadata(subtest.resourceId, id, subtest.dereferencingType)

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)
			fmt.Println(err)
			if err == nil {
				require.EqualValues(t, subtest.expectedContentStream, dereferencingResult.ContentStream)
				require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
				require.EqualValues(t, expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
				require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
				require.Empty(t, dereferencingResult.DereferencingMetadata.ResolutionError)
			} else {
				require.EqualValues(t, subtest.expectedError.Message, err.Message)
				require.EqualValues(t, subtest.expectedError.Code, err.Code)
			}
		})
	}
}

func TestDereferenceCollectionResources(t *testing.T) {
	validResource := ValidResource()
	subtests := getSubtest(types.NewDereferencedResourceList(ValidDid, []*resource.Metadata{validResource.Metadata}))

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			if !utils.IsValidResourceId(subtest.resourceId) {
				return
			}
			resourceService := services.NewResourceService(ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        ValidDid,
					MethodSpecificId: ValidIdentifier,
					Method:           ValidMethod,
				}
			}
			expectedContentType := subtest.expectedContentType
			if expectedContentType == "" {
				expectedContentType = subtest.dereferencingType
			}
			dereferencingResult, err := resourceService.DereferenceCollectionResources(id, subtest.dereferencingType)

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)
			if err == nil {
				require.EqualValues(t, subtest.expectedContentStream, dereferencingResult.ContentStream)
				require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
				require.EqualValues(t, expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
				require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
				require.Empty(t, dereferencingResult.DereferencingMetadata.ResolutionError)
			} else {
				require.EqualValues(t, subtest.expectedError.Message, err.Message)
				require.EqualValues(t, subtest.expectedError.Code, err.Code)
			}
		})
	}
}

func TestDereferenceResourceData(t *testing.T) {
	validResource := ValidResource()
	validResourceData := types.DereferencedResourceData(validResource.Resource.Data)
	subtests := getSubtest(&validResourceData)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			resourceService := services.NewResourceService(ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        ValidDid,
					MethodSpecificId: ValidIdentifier,
					Method:           ValidMethod,
				}
			}
			expectedContentType := validResource.Metadata.MediaType
			dereferencingResult, err := resourceService.DereferenceResourceData(subtest.resourceId, id, subtest.dereferencingType)

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)
			if err == nil {
				require.EqualValues(t, subtest.expectedContentStream, dereferencingResult.ContentStream)
				require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
				require.EqualValues(t, expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
				require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
				require.Empty(t, dereferencingResult.DereferencingMetadata.ResolutionError)
			} else {
				require.EqualValues(t, subtest.expectedError.Message, err.Message)
				require.EqualValues(t, subtest.expectedError.Code, err.Code)
			}
		})
	}
}
