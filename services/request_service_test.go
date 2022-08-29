package services

import (
	"fmt"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/stretchr/testify/require"
)

type MockLedgerService struct {
	Did      cheqd.Did
	Metadata cheqd.Metadata
	Resource resource.Resource
}

func NewMockLedgerService(did cheqd.Did, metadata cheqd.Metadata, resource resource.Resource) MockLedgerService {
	return MockLedgerService{
		Did:      did,
		Metadata: metadata,
		Resource: resource,
	}
}

func (ls MockLedgerService) QueryDIDDoc(did string) (cheqd.Did, cheqd.Metadata, bool, error) {
	isFound := true
	if ls.Did.Id != did {
		isFound = false
	}
	return ls.Did, ls.Metadata, isFound, nil
}

func (ls MockLedgerService) QueryResource(did string, resourceId string) (*resource.Resource, types.ErrorType) {
	if ls.Resource.Header == nil || ls.Resource.Header.Id != resourceId {
		return &resource.Resource{}, types.NotFoundError
	}
	return &ls.Resource, ""
}

func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resource.ResourceHeader, types.ErrorType) {
	if ls.Metadata.Resources == nil {
		return []*resource.ResourceHeader{}, types.NotFoundError
	}
	return []*resource.ResourceHeader{ls.Resource.Header}, ""
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}

func TestResolve(t *testing.T) {
	validDIDDoc := utils.ValidDIDDoc()
	validMetadata := utils.ValidMetadata()
	validResource := utils.ValidResource()
	subtests := []struct {
		name                   string
		ledgerService          MockLedgerService
		resolutionType         types.ContentType
		identifier             string
		method                 string
		namespace              string
		expectedDID            *types.DidDoc
		expectedMetadata       types.ResolutionDidDocMetadata
		expectedResolutionType types.ContentType
		expectedError          types.ErrorType
	}{
		{
			name:             "successful resolution",
			ledgerService:    NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      types.NewDidDoc(validDIDDoc),
			expectedMetadata: types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.ResourceHeader{validResource.Header}),
			expectedError:    "",
		},
		{
			name:             "DID not found",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NotFoundError,
		},
		{
			name:             "invalid DID",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.InvalidDIDError,
		},
		{
			name:             "invalid method",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           "not_supported_method",
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.MethodNotSupportedError,
		},
		{
			name:             "invalid namespace",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        "invalid_namespace",
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.InvalidDIDError,
		},
		{
			name:                   "representation is not supported",
			ledgerService:          NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:         "text/html,application/xhtml+xml",
			identifier:             utils.ValidIdentifier,
			method:                 utils.ValidMethod,
			namespace:              utils.ValidNamespace,
			expectedDID:            nil,
			expectedMetadata:       types.ResolutionDidDocMetadata{},
			expectedResolutionType: types.JSON,
			expectedError:          types.RepresentationNotSupportedError,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			requestService := NewRequestService("cheqd", subtest.ledgerService)
			id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier
			expectedDIDProperties := types.DidProperties{
				DidString:        id,
				MethodSpecificId: subtest.identifier,
				Method:           subtest.method,
			}
			if (subtest.resolutionType == "" || subtest.resolutionType == types.DIDJSONLD) && subtest.expectedError == "" {
				subtest.expectedDID.Context = []string{types.DIDSchemaJSONLD}
			} else if subtest.expectedDID != nil {
				subtest.expectedDID.Context = nil
			}
			expectedContentType := subtest.expectedResolutionType
			if expectedContentType == "" {
				expectedContentType = subtest.resolutionType
			}
			resolutionResult := requestService.Resolve(id, types.ResolutionOption{Accept: subtest.resolutionType})

			require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
			require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
			require.EqualValues(t, expectedContentType, resolutionResult.ResolutionMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, resolutionResult.ResolutionMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
		})
	}
}

func TestDereferencing(t *testing.T) {
	validDIDDoc := utils.ValidDIDDoc()
	validVerificationMethod := utils.ValidVerificationMethod()
	validService := utils.ValidService()
	validResource := utils.ValidResource()
	validResourceData := types.DereferencedResourceData(validResource.Data)
	validMetadata := utils.ValidMetadata()
	validFragmentMetadata := types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.ResourceHeader{})
	subtests := []struct {
		name                  string
		ledgerService         MockLedgerService
		dereferencingType     types.ContentType
		didUrl                string
		expectedContentStream types.ContentStreamI
		expectedContentType   types.ContentType
		expectedMetadata      types.ResolutionDidDocMetadata
		expectedError         types.ErrorType
	}{
		{
			name:                  "successful resolution",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			didUrl:                utils.ValidDid,
			expectedContentStream: types.NewDidDoc(validDIDDoc),
			expectedMetadata:      types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.ResourceHeader{validResource.Header}),
			expectedError:         "",
		},
		{
			name:                  "successful Secondary dereferencing (key)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			didUrl:                validVerificationMethod.Id,
			expectedContentStream: types.NewVerificationMethod(&validVerificationMethod),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         "",
		},
		{
			name:                  "successful Secondary dereferencing (service)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			didUrl:                validService.Id,
			expectedContentStream: types.NewService(&validService),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         "",
		},
		{
			name:                  "successful Primary dereferencing (resource header)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			didUrl:                utils.ValidDid + types.RESOURCE_PATH + utils.ValidResourceId + "/metadata",
			expectedContentStream: types.NewDereferencedResourceList(utils.ValidDid, []*resource.ResourceHeader{validResource.Header}),
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         "",
		},
		{
			name:                  "successful Primary dereferencing (resource list)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			didUrl:                utils.ValidDid + types.RESOURCE_PATH + "all",
			expectedContentStream: types.NewDereferencedResourceList(utils.ValidDid, []*resource.ResourceHeader{validResource.Header}),
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         "",
		},
		{
			name:                  "successful Primary dereferencing (resource data)",
			ledgerService:         NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSONLD,
			didUrl:                utils.ValidDid + types.RESOURCE_PATH + utils.ValidResourceId,
			expectedContentStream: &validResourceData,
			expectedContentType:   types.ContentType(validResource.Header.MediaType),
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         "",
		},
		{
			name:              "invalid URL",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			didUrl:            "unutils.Valid_url",
			dereferencingType: types.DIDJSONLD,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.InvalidDIDUrlError,
		},
		{
			name:              "not supported path",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            utils.ValidDid + "/unknown_path",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.RepresentationNotSupportedError,
		},
		{
			name:              "not supported query",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            utils.ValidDid + "?unknown_query",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.RepresentationNotSupportedError,
		},
		{
			name:              "key not found",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            utils.ValidDid + "#notFoundKey",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NotFoundError,
		},
		{
			name:              "resource not found",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            utils.ValidDid + types.RESOURCE_PATH + "00000000-0000-0000-0000-000000000000",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NotFoundError,
		},
		{
			name: 				"successful primary service dereferencing",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"Primary service dereferencing(not found)",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=serv",
			expectedMetadata:  	types.ResolutionDidDocMetadata{},
			expectedError:     	types.NotFoundError,
		},
		{
			name: 				"Primary service dereferencing(hash simpol)",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1#flag",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint + "#flag"),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"Primary service dereferencing(relativeRef)",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1&relativeRef=/some/path?some_query",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint + "/some/path?some_query"),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"Primary dereferencing(relativeRef + hash flag)",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1&relativeRef=/some/path?some_query#flag",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint + "/some/path?some_query#flag"),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			requestService := NewRequestService("cheqd", subtest.ledgerService)
			var expectedDIDProperties types.DidProperties
			if subtest.expectedError != types.InvalidDIDUrlError {
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

			fmt.Println(" dereferencingResult   " + subtest.didUrl)

			dereferencingResult := requestService.Dereference(subtest.didUrl, types.DereferencingOption{Accept: subtest.dereferencingType})

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)

			require.EqualValues(t, subtest.expectedContentStream, dereferencingResult.ContentStream)
			require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
			require.EqualValues(t, expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, dereferencingResult.DereferencingMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
			require.EqualValues(t, subtest.expectedError.GetStatusCode(), dereferencingResult.DereferencingMetadata.ResolutionError.GetStatusCode())
		})
	}
}
