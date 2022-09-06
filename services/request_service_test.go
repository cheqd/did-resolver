package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestResolveDIDDoc(t *testing.T) {
	validDIDDoc := utils.ValidDIDDoc()
	validMetadata := utils.ValidMetadata()
	validResource := utils.ValidResource()
	validDIDResolution := types.NewDidDoc(validDIDDoc)
	subtests := []struct {
		name                   string
		ledgerService          utils.MockLedgerService
		resolutionType         types.ContentType
		did                    string
		expectedDID            *types.DidDoc
		expectedMetadata       types.ResolutionDidDocMetadata
		expectedResolutionType types.ContentType
		expectedError          *types.IdentityError
	}{
		{
			name:             "successful resolution",
			ledgerService:    utils.NewMockLedgerService(validDIDDoc, validMetadata, validResource),
			resolutionType:   types.DIDJSONLD,
			did:              utils.ValidDid,
			expectedDID:      &validDIDResolution,
			expectedMetadata: types.NewResolutionDidDocMetadata(utils.ValidDid, validMetadata, []*resource.ResourceHeader{validResource.Header}),
			expectedError:    nil,
		},
		{
			name:             "DID not found",
			ledgerService:    utils.NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}, resource.Resource{}),
			resolutionType:   types.DIDJSONLD,
			did:              utils.ValidDid,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(utils.ValidDid, types.DIDJSONLD, nil, false),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			
			// Setup
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			context := e.NewContext(req, rec)
			context.SetPath("/1.0/identifiers/:did")
			context.SetParamNames("did")
			context.SetParamValues(subtest.did)

			requestService := NewRequestService("cheqd", subtest.ledgerService)
			expectedDIDProperties := types.DidProperties{
				DidString:        utils.ValidDid,
				MethodSpecificId: utils.ValidIdentifier,
				Method:           utils.ValidMethod,
			}
			if (subtest.resolutionType == "" || subtest.resolutionType == types.DIDJSONLD) && subtest.expectedError == nil {
				subtest.expectedDID.Context = []string{types.DIDSchemaJSONLD}
			} else if subtest.expectedDID != nil {
				subtest.expectedDID.Context = nil
			}
			expectedContentType := subtest.expectedResolutionType
			if expectedContentType == "" {
				expectedContentType = subtest.resolutionType
			}
			err := requestService.ResolveDIDDoc(context)
			var resolutionResult types.DidResolution
			json.Unmarshal(rec.Body.Bytes(), &resolutionResult)

			if subtest.expectedError != nil {
				require.EqualValues(t, subtest.expectedError, err)
			} else {
				require.Empty(t, err)
				require.EqualValues(t, subtest.expectedError, err)
				require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
				require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
				require.EqualValues(t, expectedContentType, resolutionResult.ResolutionMetadata.ContentType)
				require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
			}
		})
	}
}
