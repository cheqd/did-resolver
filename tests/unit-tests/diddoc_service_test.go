package tests

import (
	"fmt"
	"net/url"
	"testing"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
)

func TestDIDDocFragment(t *testing.T) {
	validDIDDoc := ValidDIDDoc()
	DIDDoc := types.NewDidDoc(&validDIDDoc)

	subtests := []struct {
		name             string
		fragmentId       string
		didDoc           types.DidDoc
		expectedFragment types.ContentStreamI
	}{
		{
			name:             "successful VerificationMethod finding",
			fragmentId:       DIDDoc.VerificationMethod[0].Id,
			didDoc:           DIDDoc,
			expectedFragment: &DIDDoc.VerificationMethod[0],
		},
		{
			name:             "successful Service finding",
			fragmentId:       DIDDoc.Service[0].Id,
			didDoc:           DIDDoc,
			expectedFragment: &DIDDoc.Service[0],
		},
		{
			name:             "Fragment is not found",
			fragmentId:       "fake_id",
			didDoc:           DIDDoc,
			expectedFragment: nil,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			didDocService := services.DIDDocService{}

			fragment := didDocService.GetDIDFragment(subtest.fragmentId, subtest.didDoc)

			require.EqualValues(t, subtest.expectedFragment, fragment)
		})
	}
}

func TestResolve(t *testing.T) {
	validDIDDoc := ValidDIDDoc()
	validDIDDocResolution := types.NewDidDoc(&validDIDDoc)
	validMetadata := ValidMetadata()
	validResource := ValidResource()
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
		expectedError          *types.IdentityError
	}{
		{
			name:             "successful resolution",
			ledgerService:    NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      &validDIDDocResolution,
			expectedMetadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
		{
			name:             "DID not found",
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid DID",
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid method",
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           "not_supported_method",
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:             "invalid namespace",
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        "invalid_namespace",
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:                   "representation is not supported",
			ledgerService:          NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:         "text/html,application/xhtml+xml",
			identifier:             ValidIdentifier,
			method:                 ValidMethod,
			namespace:              ValidNamespace,
			expectedDID:            nil,
			expectedMetadata:       types.ResolutionDidDocMetadata{},
			expectedResolutionType: types.JSON,
			expectedError:          types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, false),
		},
	}

	for _, subtest := range subtests {
		fmt.Printf("Testing %s", subtest.name)
		id := "did:" + subtest.method + ":" + subtest.namespace + ":" + subtest.identifier
		t.Run(subtest.name, func(t *testing.T) {
			diddocService := services.NewDIDDocService("cheqd", subtest.ledgerService)
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
	validDIDDoc := ValidDIDDoc()
	validVerificationMethod := ValidVerificationMethod()
	validService := ValidService()
	validResource := ValidResource()
	validMetadata := ValidMetadata()
	validFragmentMetadata := types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{})
	validQuery, _ := url.ParseQuery("attr=value")
	subtests := []struct {
		name                  string
		ledgerService         MockLedgerService
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
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			did:                   ValidDid,
			fragmentId:            validVerificationMethod.Id,
			expectedContentStream: types.NewVerificationMethod(&validVerificationMethod),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
		{
			name:                  "successful Secondary dereferencing (service)",
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			did:                   ValidDid,
			fragmentId:            validService.Id,
			expectedContentStream: types.NewService(&validService),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
		{
			name:                  "not supported query",
			ledgerService:         NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   ValidDid,
			queries:               validQuery,
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, false),
		},
		{
			name:                  "key not found",
			ledgerService:         NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   ValidDid,
			fragmentId:            "notFoundKey",
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			diddocService := services.NewDIDDocService("cheqd", subtest.ledgerService)
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
