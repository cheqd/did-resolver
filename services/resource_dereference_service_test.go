package services

import (
	"fmt"
	"strings"
	"testing"

	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	name                  string
	ledgerService         utils.MockLedgerService
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
	validDIDDoc := utils.ValidDIDDoc()
	validResource := utils.ValidResource()
	validMetadata := utils.ValidMetadata()
	return []TestCase{
		{
			name:                  "successful dereferencing for resource",
			ledgerService:         utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            utils.ValidIdentifier,
			method:                utils.ValidMethod,
			namespace:             utils.ValidNamespace,
			resourceId:            utils.ValidResourceId,
			expectedContentStream: validContentStream,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         nil,
		},
		{
			name:                  "successful dereferencing for resource (upper case UUID)",
			ledgerService:         utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			identifier:            utils.ValidIdentifier,
			method:                utils.ValidMethod,
			namespace:             utils.ValidNamespace,
			resourceId:            strings.ToUpper(utils.ValidResourceId),
			expectedContentStream: validContentStream,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         nil,
		},
		{
			name:              "resource not found",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        utils.ValidIdentifier,
			method:            utils.ValidMethod,
			namespace:         utils.ValidNamespace,
			resourceId:        "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewNotFoundError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid resource id",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        utils.ValidIdentifier,
			method:            utils.ValidMethod,
			namespace:         utils.ValidNamespace,
			resourceId:        "invalid-resource-id",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid type",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.JSON,
			identifier:        utils.ValidIdentifier,
			method:            utils.ValidMethod,
			namespace:         utils.ValidNamespace,
			resourceId:        utils.ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewRepresentationNotSupportedError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid namespace",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        utils.ValidIdentifier,
			method:            utils.ValidMethod,
			namespace:         "invalid-namespace",
			resourceId:        utils.ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid method",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        utils.ValidIdentifier,
			method:            "invalid-method",
			namespace:         utils.ValidNamespace,
			resourceId:        utils.ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
		{
			name:              "invalid identifier",
			ledgerService:     utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType: types.DIDJSONLD,
			identifier:        "invalid-identifier",
			method:            utils.ValidMethod,
			namespace:         utils.ValidNamespace,
			resourceId:        utils.ValidResourceId,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NewInvalidDIDUrlError(utils.ValidDid, types.DIDJSONLD, nil, true),
		},
	}
}

func TestDereferenceResourceMetadata(t *testing.T) {
	validResource := utils.ValidResource()
	subtests := getSubtest(types.NewDereferencedResourceList(utils.ValidDid, []*resource.Metadata{validResource.Metadata}))

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			resourceService := NewResourceService(utils.ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        utils.ValidDid,
					MethodSpecificId: utils.ValidIdentifier,
					Method:           utils.ValidMethod,
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
	validResource := utils.ValidResource()
	subtests := getSubtest(types.NewDereferencedResourceList(utils.ValidDid, []*resource.Metadata{validResource.Metadata}))

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			if !utils.IsValidResourceId(subtest.resourceId) {
				return
			}
			resourceService := NewResourceService(utils.ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        utils.ValidDid,
					MethodSpecificId: utils.ValidIdentifier,
					Method:           utils.ValidMethod,
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
	validResource := utils.ValidResource()
	validResourceData := types.DereferencedResourceData(validResource.Resource.Data)
	subtests := getSubtest(&validResourceData)

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			resourceService := NewResourceService(utils.ValidMethod, subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier

			var expectedDIDProperties types.DidProperties
			if subtest.expectedError == nil {
				expectedDIDProperties = types.DidProperties{
					DidString:        utils.ValidDid,
					MethodSpecificId: utils.ValidIdentifier,
					Method:           utils.ValidMethod,
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
