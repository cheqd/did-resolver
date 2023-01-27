package services

import (
	"fmt"
	"net/url"
	"testing"

	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/stretchr/testify/require"
)

func TestDIDDocFragment(t *testing.T) {
	validDIDDoc := types.NewDidDoc(utils.ValidDIDDoc())

	subtests := []struct {
		name             string
		fragmentId       string
		didDoc           types.DidDoc
		expectedFragment types.ContentStreamI
	}{
		{
			name:             "successful VerificationMethod finding",
			fragmentId:       validDIDDoc.VerificationMethod[0].Id,
			didDoc:           validDIDDoc,
			expectedFragment: &validDIDDoc.VerificationMethod[0],
		},
		{
			name:             "successful Service finding",
			fragmentId:       validDIDDoc.Service[0].Id,
			didDoc:           validDIDDoc,
			expectedFragment: &validDIDDoc.Service[0],
		},
		{
			name:             "Fragment is not found",
			fragmentId:       "fake_id",
			didDoc:           validDIDDoc,
			expectedFragment: nil,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			didDocService := DIDDocService{}

			fragment := didDocService.GetDIDFragment(subtest.fragmentId, subtest.didDoc)

			require.EqualValues(t, subtest.expectedFragment, fragment)
		})
	}
}

func TestResolve(t *testing.T) {
	validDIDDoc := utils.ValidDIDDoc()
	validDIDDocResolution := types.NewDidDoc(validDIDDoc)
	validMetadata := utils.ValidMetadata()
	validResource := utils.ValidResource()
	subtests := []struct {
		name                   string
		ledgerService          utils.MockLedgerService
		resolutionType         types.ContentType
		identifier             string
		method                 string
		namespace              string
		expectedDID            *types.DidDoc
		expectedMetadata       types.ResolutionDidDocMetadata
		expectedResolutionType types.ContentType
		expectedError          *types.IdentityError
	}{
		{
			name:             "successful resolution",
			ledgerService:    utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      &validDIDDocResolution,
			expectedMetadata: types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
		{
			name:             "DID not found",
			ledgerService:    utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid DID",
			ledgerService:    utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           utils.ValidMethod,
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewInvalidDIDError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid method",
			ledgerService:    utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           "not_supported_method",
			namespace:        utils.ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewMethodNotSupportedError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid namespace",
			ledgerService:    utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       utils.ValidIdentifier,
			method:           utils.ValidMethod,
			namespace:        "invalid_namespace",
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewInvalidDIDError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:                   "representation is not supported",
			ledgerService:          utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:         "text/html,application/xhtml+xml",
			identifier:             utils.ValidIdentifier,
			method:                 utils.ValidMethod,
			namespace:              utils.ValidNamespace,
			expectedDID:            nil,
			expectedMetadata:       types.ResolutionDidDocMetadata{},
			expectedResolutionType: types.JSON,
			expectedError:          types.NewRepresentationNotSupportedError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
	}

	for _, subtest := range subtests {
		fmt.Printf("Testing %s", subtest.name)
		id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier
		t.Run(subtest.name, func(t *testing.T) {
			diddocService := NewDIDDocService("cheqd", subtest.ledgerService)
			expectedDIDProperties := types.DidProperties{
				DidString:        id,
				MethodSpecificId: subtest.identifier,
				Method:           subtest.method,
			}
			if (subtest.resolutionType == "" || subtest.resolutionType == types.DIDJSONLD) && subtest.expectedError == nil {
				subtest.expectedDID.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
			} else if subtest.expectedDID != nil {
				subtest.expectedDID.Context = nil
			}
			expectedContentType := subtest.expectedResolutionType
			if expectedContentType == "" {
				expectedContentType = subtest.resolutionType
			}
			resolutionResult, err := diddocService.Resolve(id, "", subtest.resolutionType)
			// print(resolutionResult.Did.Id)
			if subtest.expectedError != nil {
				require.EqualValues(t, subtest.expectedError.Code, err.Code)
				require.EqualValues(t, subtest.expectedError.Message, err.Message)
			} else {
				require.Empty(t, err)
				require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
				require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
				require.EqualValues(t, expectedContentType, resolutionResult.ResolutionMetadata.ContentType)
				require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
			}
		})
	}
}

func TestDereferencing(t *testing.T) {
	validDIDDoc := utils.ValidDIDDoc()
	validVerificationMethod := utils.ValidVerificationMethod()
	validService := utils.ValidService()
	validResource := utils.ValidResource()
	validMetadata := utils.ValidMetadata()
	validFragmentMetadata := types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.Metadata{})
	validQuery, _ := url.ParseQuery("attr=value")
	subtests := []struct {
		name                  string
		ledgerService         utils.MockLedgerService
		dereferencingType     types.ContentType
		did                   string
		fragmentId            string
		queries               url.Values
		expectedContentStream types.ContentStreamI
		expectedMetadata      types.ResolutionDidDocMetadata
		expectedContentType   types.ContentType
		expectedError         *types.IdentityError
	}{
		{
			name:                  "successful Secondary dereferencing (key)",
			ledgerService:         utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			did:                   utils.ValidDid,
			fragmentId:            validVerificationMethod.Id,
			expectedContentStream: types.NewVerificationMethod(&validVerificationMethod),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
		{
			name:                  "successful Secondary dereferencing (service)",
			ledgerService:         utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			dereferencingType:     types.DIDJSON,
			did:                   utils.ValidDid,
			fragmentId:            validService.Id,
			expectedContentStream: types.NewService(&validService),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
		{
			name:                  "not supported query",
			ledgerService:         utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   utils.ValidDid,
			queries:               validQuery,
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewRepresentationNotSupportedError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:                  "key not found",
			ledgerService:         utils.NewMockLedgerService(did.DidDoc{}, did.Metadata{}, resource.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   utils.ValidDid,
			fragmentId:            "notFoundKey",
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewNotFoundError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			diddocService := NewDIDDocService("cheqd", subtest.ledgerService)
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

			result, err := diddocService.ProcessDIDRequest(subtest.did, subtest.fragmentId, subtest.queries, nil, subtest.dereferencingType)
			dereferencingResult, _ := result.(*types.DidDereferencing)

			fmt.Println(subtest.name + ": dereferencingResult:")
			fmt.Println(dereferencingResult)

			if subtest.expectedError != nil {
				require.EqualValues(t, subtest.expectedError.Code, err.Code)
				require.EqualValues(t, subtest.expectedError.Message, err.Message)
			} else {
				require.Empty(t, err)
				require.EqualValues(t, subtest.expectedContentStream, dereferencingResult.ContentStream)
				require.EqualValues(t, subtest.expectedMetadata, dereferencingResult.Metadata)
				require.EqualValues(t, expectedContentType, dereferencingResult.DereferencingMetadata.ContentType)
				require.Empty(t, dereferencingResult.DereferencingMetadata.ResolutionError)
				require.EqualValues(t, expectedDIDProperties, dereferencingResult.DereferencingMetadata.DidProperties)
			}
		})
	}
}
