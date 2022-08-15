package services

import (
	"errors"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/stretchr/testify/require"
)

func TestQueryDIDDoc(t *testing.T) {
	subtests := []struct {
		name             string
		did              string
		expectedDID      cheqd.Did
		expectedMetadata cheqd.Metadata
		expectedIsFound  bool
		expectedError    error
	}{
		{
			name:             "DeadlineExceeded",
			did:              "fake did",
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedIsFound:  false,
			expectedError:    errors.New("namespace not supported: "),
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ledgerService := NewLedgerService()
			didDoc, metadata, isFound, err := ledgerService.QueryDIDDoc("fake did")
			require.EqualValues(t, subtest.expectedDID, didDoc)
			require.EqualValues(t, subtest.expectedMetadata, metadata)
			require.EqualValues(t, subtest.expectedIsFound, isFound)
			require.EqualValues(t, subtest.expectedError, err)
		})
	}
}

func TestQueryResource(t *testing.T) {
	subtests := []struct {
		name             string
		collectionDid    string
		resourceId       string
		expectedResource resource.Resource
		expectedError    types.ErrorType
	}{
		{
			name:             "DeadlineExceeded",
			collectionDid:    "321",
			resourceId:       "123",
			expectedResource: resource.Resource{},
			expectedError:    types.InvalidDIDError,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			ledgerService := NewLedgerService()
			resource, errorType := ledgerService.QueryResource(subtest.collectionDid, subtest.resourceId)
			require.EqualValues(t, &subtest.expectedResource, resource)
			require.EqualValues(t, subtest.expectedError, errorType)
		})
	}
}
