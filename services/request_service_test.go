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

func (ls MockLedgerService) QueryResource(collectionDid string, resourceId string) (resource.Resource, bool, error) {
	isFound := true
	if ls.Resource.Header == nil {
		isFound = false
	}
	return ls.Resource, isFound, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}

func TestResolve(t *testing.T) {
	validDIDDoc := validDIDDoc()
	subtests := []struct {
		name             string
		ledgerService    MockLedgerService
		resolutionType   types.ContentType
		identifier       string
		method           string
		namespace        string
		expectedDID      cheqd.Did
		expectedMetadata cheqd.Metadata
		expectedError    types.ErrorType
	}{
		{
			name:             "successful resolution",
			ledgerService:    NewMockLedgerService(validDIDDoc, validMetadata(), resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			namespace:        validNamespace,
			expectedDID:      validDIDDoc,
			expectedMetadata: validMetadata(),
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
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionNotFound,
		},
		{
			name:             "invalid DID",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           validMethod,
			namespace:        validNamespace,
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionInvalidDID,
		},
		{
			name:             "invalid method",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           "not_supported_method",
			namespace:        validNamespace,
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionMethodNotSupported,
		},
		{
			name:             "invalid namespace",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			namespace:        "invalid_namespace",
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionInvalidDID,
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
			if (subtest.resolutionType == types.DIDJSONLD || subtest.resolutionType == types.JSONLD) && subtest.expectedError == "" {
				subtest.expectedDID.Context = []string{types.DIDSchemaJSONLD}
			}

			resolutionResult, err := requestService.Resolve(id, types.ResolutionOption{Accept: subtest.resolutionType})

			fmt.Println(subtest.name + ": resolutionResult:")
			fmt.Println(resolutionResult.Did.VerificationMethod)
			fmt.Println(subtest.expectedDID.VerificationMethod)
			require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
			require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
			require.EqualValues(t, subtest.resolutionType, resolutionResult.ResolutionMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, resolutionResult.ResolutionMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
			require.Empty(t, err)
		})
	}
}

func TestDereferencing(t *testing.T) {
	validDIDDoc := validDIDDoc()
	validVerificationMethod := validVerificationMethod()
	validService := validService()
	validResource := validResource()
	validChecksum, _ := json.Marshal(validResource.Header.Checksum)
	validData, _ := json.Marshal(validResource.Data)
	subtests := []struct {
		name                  string
		ledgerService         MockLedgerService
		dereferencingType     types.ContentType
		didUrl                string
		expectedContentStream string
		expectedMetadata      cheqd.Metadata
		expectedError         types.ErrorType
	}{
		{
			name:              "successful resolution",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata(), validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid,
			expectedContentStream: fmt.Sprintf("{\"@context\":[\"%s\"],\"id\":\"%s\",\"verificationMethod\":[{\"id\":\"%s\",\"type\":\"%s\",\"controller\":\"%s\",\"publicKeyJwk\":%s}],\"service\":[{\"id\":\"%s\",\"type\":\"%s\",\"serviceEndpoint\":\"%s\"}]}",
				types.DIDSchemaJSONLD, validDid, validVerificationMethod.Id, validVerificationMethod.Type, validVerificationMethod.Controller, validPubKeyJWK, validService.Id, validService.Type, validService.ServiceEndpoint),
			expectedMetadata: validMetadata(),
			expectedError:    "",
		},
		{
			name:              "successful Secondary dereferencing (key)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata(), validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validVerificationMethod.Id,
			expectedContentStream: fmt.Sprintf("{\"@context\":\"%s\",\"id\":\"%s\",\"type\":\"%s\",\"controller\":\"%s\",\"publicKeyJwk\":%s}",
				types.DIDSchemaJSONLD, validVerificationMethod.Id, validVerificationMethod.Type, validVerificationMethod.Controller, validPubKeyJWK),
			expectedMetadata: validMetadata(),
			expectedError:    "",
		},
		{
			name:              "successful Secondary dereferencing (service)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata(), validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validService.Id,
			expectedContentStream: fmt.Sprintf("{\"@context\":\"%s\",\"id\":\"%s\",\"type\":\"%s\",\"serviceEndpoint\":\"%s\"}",
				types.DIDSchemaJSONLD, validService.Id, validService.Type, validService.ServiceEndpoint),
			expectedMetadata: validMetadata(),
			expectedError:    "",
		},
		{
			name:              "successful Primary dereferencing (resource)",
			ledgerService:     NewMockLedgerService(validDIDDoc, validMetadata(), validResource),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "/resource/" + validResourceId,
			expectedContentStream: fmt.Sprintf("{\"@context\":[\"%s\"],\"collectionId\":\"%s\",\"id\":\"%s\",\"name\":\"%s\",\"resourceType\":\"%s\",\"mediaType\":\"%s\",\"checksum\":%s,\"data\":%s}",
				types.DIDSchemaJSONLD, validResource.Header.CollectionId, validResource.Header.Id, validResource.Header.Name, validResource.Header.ResourceType, validResource.Header.MediaType, validChecksum, validData),
			expectedMetadata: cheqd.Metadata{},
			expectedError:    "",
		},
		{
			name:              "invalid URL",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			didUrl:            "unvalid_url",
			dereferencingType: types.DIDJSONLD,
			expectedMetadata:  cheqd.Metadata{},
			expectedError:     types.DereferencingInvalidDIDUrl,
		},
		{
			name:              "not supported path",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "/unknown_path",
			expectedMetadata:  cheqd.Metadata{},
			expectedError:     types.DereferencingNotSupported,
		},
		{
			name:              "not supported query",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "?unknown_query",
			expectedMetadata:  cheqd.Metadata{},
			expectedError:     types.DereferencingNotSupported,
		},
		{
			name:              "key not found",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "#notFoundKey",
			expectedMetadata:  cheqd.Metadata{},
			expectedError:     types.DereferencingNotFound,
		},
		{
			name:              "resource not found",
			ledgerService:     NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			dereferencingType: types.DIDJSONLD,
			didUrl:            validDid + "/resource/00000000-0000-0000-0000-000000000000",
			expectedMetadata:  cheqd.Metadata{},
			expectedError:     types.DereferencingNotFound,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			requestService := NewRequestService("cheqd", subtest.ledgerService)
			var expectedDIDProperties types.DidProperties
			if subtest.expectedError != types.DereferencingInvalidDIDUrl {
				expectedDIDProperties = types.DidProperties{
					DidString:        validDid,
					MethodSpecificId: validIdentifier,
					Method:           validMethod,
				}
			}

			fmt.Println(" dereferencingResult   " + subtest.didUrl)

			dereferencingResult, err := requestService.Dereference(subtest.didUrl, types.DereferencingOption{Accept: subtest.dereferencingType})

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)
			require.EqualValues(t, string(subtest.expectedContentStream), string(dereferencingResult.ContentStream))
			require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
			require.EqualValues(t, subtest.dereferencingType, dereferencingResult.DereferencingMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, dereferencingResult.DereferencingMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
			require.Empty(t, err)
		})
	}
}
