package services

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
)

const (
	validIdentifier = "N22KY2Dyvmuu2Pyy"
	validMethod     = "cheqd"
	validNamespace  = "mainnet"
	validDid        = "did:" + validMethod + ":" + validNamespace + ":" + validIdentifier
	validResourceId = "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
	validPubKeyJWK  = "{" +
		"\"crv\":\"Ed25519\"," +
		"\"kid\":\"_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A\"," +
		"\"kty\":\"OKP\"," +
		"\"x\":\"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"" +
		"}"
)

func validVerificationMethod() cheqd.VerificationMethod {
	return cheqd.VerificationMethod{
		Id:           validDid + "#key-1",
		Type:         "JsonWebKey2020",
		Controller:   validDid,
		PublicKeyJwk: cheqd.JSONToPubKeyJWK(validPubKeyJWK),
	}
}

func validService() cheqd.Service {
	return cheqd.Service{
		Id:              validDid + "#service-1",
		Type:            "DIDCommMessaging",
		ServiceEndpoint: "endpoint",
	}
}

func validDIDDoc() cheqd.Did {
	service := validService()
	verificationMethod := validVerificationMethod()

	return cheqd.Did{
		Id:                 validDid,
		VerificationMethod: []*cheqd.VerificationMethod{&verificationMethod},
		Service:            []*cheqd.Service{&service},
	}
}

func validResource() resource.Resource {
	data := []byte("{\"attr\":[\"name\",\"age\"]}")
	return resource.Resource{
		Header: &resource.ResourceHeader{
			CollectionId: validIdentifier,
			Id:           validResourceId,
			Name:         "Existing Resource Name",
			ResourceType: "CL-Schema",
			MediaType:    "application/json",
			Checksum:     sha256.New().Sum(data),
		},
		Data: data,
	}
}

func validMetadata() cheqd.Metadata {
	return cheqd.Metadata{VersionId: "test_version_id", Deactivated: false, Resources: []string{validResourceId}}
}

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

func (ls MockLedgerService) QueryResource(did string, resourceId string) (resource.Resource, bool, error) {
	isFound := true
	if ls.Resource.Header == nil {
		isFound = false
	}
	return ls.Resource, isFound, nil
}

func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resource.ResourceHeader, error) {
	if ls.Metadata.Resources == nil {
		return []*resource.ResourceHeader{}, nil
	}
	return []*resource.ResourceHeader{ls.Resource.Header}, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}

func TestResolve(t *testing.T) {
	validDIDDoc := validDIDDoc()
	validMetadata := validMetadata()
	validResource := validResource()
	subtests := []struct {
		name                   string
		ledgerService          MockLedgerService
		resolutionType         types.ContentType
		identifier             string
		method                 string
		namespace              string
		expectedDID            cheqd.Did
		expectedMetadata       types.ResolutionDidDocMetadata
		expectedResolutionType types.ContentType
		expectedError          types.ErrorType
	}{
		{
			name:             "successful resolution",
			ledgerService:    NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			namespace:        validNamespace,
			expectedDID:      validDIDDoc,
			expectedMetadata: types.NewResolutionDidDocMetadata(validDid, validMetadata, []*resource.ResourceHeader{validResource.Header}),
			expectedError:    "",
		},
		{
			name:             "DID not found",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			namespace:        validNamespace,
			expectedDID:      cheqd.Did{},
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NotFoundError,
		},
		{
			name:             "invalid DID",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           validMethod,
			namespace:        validNamespace,
			expectedDID:      cheqd.Did{},
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.InvalidDIDError,
		},
		{
			name:             "invalid method",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           "not_supported_method",
			namespace:        validNamespace,
			expectedDID:      cheqd.Did{},
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.MethodNotSupportedError,
		},
		{
			name:             "invalid namespace",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			namespace:        "invalid_namespace",
			expectedDID:      cheqd.Did{},
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.InvalidDIDError,
		},
		{
			name:                   "representation is not supported",
			ledgerService:          NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:         "text/html,application/xhtml+xml",
			identifier:             validIdentifier,
			method:                 validMethod,
			namespace:              validNamespace,
			expectedDID:            cheqd.Did{},
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
			} else {
				subtest.expectedDID.Context = nil
			}
			expectedContentType := subtest.expectedResolutionType
			if expectedContentType == "" {
				expectedContentType = subtest.resolutionType
			}
			resolutionResult := requestService.Resolve(id, types.ResolutionOption{Accept: subtest.resolutionType})

			fmt.Println(subtest.name + ": resolutionResult:")
			fmt.Println(resolutionResult.Did.VerificationMethod)
			fmt.Println(subtest.expectedDID.VerificationMethod)
			require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
			require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
			require.EqualValues(t, expectedContentType, resolutionResult.ResolutionMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, resolutionResult.ResolutionMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
		})
	}
}

func TestDereferencing(t *testing.T) {
	validDIDDoc := validDIDDoc()
	validVerificationMethod := validVerificationMethod()
	validService := validService()
	validResource := validResource()
	validChecksum := fmt.Sprintf("%x", validResource.Header.Checksum)
	validData, _ := json.Marshal(validResource.Data)
	validMetadata := validMetadata()
	validFragmentMetadata := types.NewResolutionDidDocMetadata(validDid, validMetadata, []*resource.ResourceHeader{})
	subtests := []struct {
		name                  string
		ledgerService         MockLedgerService
		dereferencingType     types.ContentType
		didUrl                string
		expectedContentStream string
		expectedMetadata      types.ResolutionDidDocMetadata
		expectedError         types.ErrorType
	}{
		{
			name:              "successful resolution",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid,
			expectedContentStream: fmt.Sprintf("{\"@context\":[\"%s\"],\"id\":\"%s\",\"verificationMethod\":[{\"id\":\"%s\",\"type\":\"%s\",\"controller\":\"%s\",\"publicKeyJwk\":%s}],\"service\":[{\"id\":\"%s\",\"type\":\"%s\",\"serviceEndpoint\":\"%s\"}]}",
				types.DIDSchemaJSONLD, validDid, validVerificationMethod.Id, validVerificationMethod.Type, validVerificationMethod.Controller, validPubKeyJWK, validService.Id, validService.Type, validService.ServiceEndpoint),
			expectedMetadata: types.NewResolutionDidDocMetadata(validDid, validMetadata, []*resource.ResourceHeader{validResource.Header}),
			expectedError:    "",
		},
		{
			name:              "successful Secondary dereferencing (key)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validVerificationMethod.Id,
			expectedContentStream: fmt.Sprintf("{\"@context\":\"%s\",\"id\":\"%s\",\"type\":\"%s\",\"controller\":\"%s\",\"publicKeyJwk\":%s}",
				types.DIDSchemaJSONLD, validVerificationMethod.Id, validVerificationMethod.Type, validVerificationMethod.Controller, validPubKeyJWK),
			expectedMetadata: validFragmentMetadata,
			expectedError:    "",
		},
		{
			name:              "successful Secondary dereferencing (service)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validService.Id,
			expectedContentStream: fmt.Sprintf("{\"@context\":\"%s\",\"id\":\"%s\",\"type\":\"%s\",\"serviceEndpoint\":\"%s\"}",
				types.DIDSchemaJSONLD, validService.Id, validService.Type, validService.ServiceEndpoint),
			expectedMetadata: validFragmentMetadata,
			expectedError:    "",
		},
		{
			name:              "successful Primary dereferencing (resource)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + types.RESOURCE_PATH + validResourceId,
			expectedContentStream: fmt.Sprintf("{\"@context\":[\"%s\"],\"collectionId\":\"%s\",\"id\":\"%s\",\"name\":\"%s\",\"resourceType\":\"%s\",\"mediaType\":\"%s\",\"checksum\":\"%s\",\"data\":%s}",
				types.DIDSchemaJSONLD, validResource.Header.CollectionId, validResource.Header.Id, validResource.Header.Name, validResource.Header.ResourceType, validResource.Header.MediaType, validChecksum, validData),
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    "",
		},
		{
			name:              "invalid URL",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			didUrl:            "unvalid_url",
			dereferencingType: types.DIDJSONLD,
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.InvalidDIDUrlError,
		},
		{
			name:              "not supported path",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "/unknown_path",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.RepresentationNotSupportedError,
		},
		{
			name:              "not supported query",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "?unknown_query",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.RepresentationNotSupportedError,
		},
		{
			name:              "key not found",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "#notFoundKey",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NotFoundError,
		},
		{
			name:              "resource not found",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + types.RESOURCE_PATH + "00000000-0000-0000-0000-000000000000",
			expectedMetadata:  types.ResolutionDidDocMetadata{},
			expectedError:     types.NotFoundError,
		},
		{
			name: 				"get clean endpoint",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"sercice not found",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=serv",
			expectedMetadata:  	types.ResolutionDidDocMetadata{},
			expectedError:     	types.NotFoundError,
		},
		{
			name: 				"hash simpol",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1#flag",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint + "#flag"),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"relativeRef",
			ledgerService: 		NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType: 	types.DIDJSONLD,
			didUrl:            	validDid + "?service=service-1&relativeRef=/some/path?some_query",
			expectedContentStream: fmt.Sprintf("\"%s\"", validService.ServiceEndpoint + "/some/path?some_query"),
			expectedMetadata: 	validFragmentMetadata,
			expectedError:    	"",
		},
		{
			name: 				"relativeRef + hash simpol",
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
					DidString:        validDid,
					MethodSpecificId: validIdentifier,
					Method:           validMethod,
				}
			}

			fmt.Println(" dereferencingResult   " + subtest.didUrl)

			dereferencingResult := requestService.Dereference(subtest.didUrl, types.DereferencingOption{Accept: subtest.dereferencingType})

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)
			require.EqualValues(t, string(subtest.expectedContentStream), string(dereferencingResult.ContentStream))
			require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
			require.EqualValues(t, subtest.dereferencingType, dereferencingResult.DereferencingMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, dereferencingResult.DereferencingMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
		})
	}
}
