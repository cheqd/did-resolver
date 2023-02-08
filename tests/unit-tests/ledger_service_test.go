package tests

import (
	"testing"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resource "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
)

func TestQueryDIDDoc(t *testing.T) {
	subtests := []struct {
		name                     string
		did                      string
		expectedDidDocWithMetada *didTypes.DidDocWithMetadata
		expectedIsFound          bool
		expectedError            error
	}{
		{
			name:                     "DeadlineExceeded",
			did:                      "fake did",
			expectedDidDocWithMetada: nil,
			expectedIsFound:          false,
			expectedError:            types.NewInvalidDIDError("fake did", types.JSON, nil, false),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ledgerService := services.NewLedgerService()
			didDocWithMetadata, err := ledgerService.QueryDIDDoc("fake did", "")
			require.EqualValues(t, subtest.expectedDidDocWithMetada, didDocWithMetadata)
			require.EqualValues(t, subtest.expectedError.Error(), err.Error())
		})
	}
}

func TestQueryResource(t *testing.T) {
	subtests := []struct {
		name             string
		collectionDid    string
		resourceId       string
		expectedResource *resource.ResourceWithMetadata
		expectedError    error
	}{
		{
			name:             "DeadlineExceeded",
			collectionDid:    "321",
			resourceId:       "123",
			expectedResource: nil,
			expectedError:    types.NewInvalidDIDError("321", types.JSON, nil, true),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ledgerService := services.NewLedgerService()
			resource, err := ledgerService.QueryResource(subtest.collectionDid, subtest.resourceId)
			require.EqualValues(t, subtest.expectedResource, resource)
			require.EqualValues(t, subtest.expectedError.Error(), err.Error())
		})
	}
}
