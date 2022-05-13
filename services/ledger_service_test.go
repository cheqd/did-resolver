package services

import (
	"errors"
	"testing"
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
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
			timeout, err := time.ParseDuration("5s")
			require.NoError(t, err)

			ledgerService := NewLedgerService(timeout)
			didDoc, metadata, isFound, err := ledgerService.QueryDIDDoc("fake did")
			require.EqualValues(t, subtest.expectedDID, didDoc)
			require.EqualValues(t, subtest.expectedMetadata, metadata)
			require.EqualValues(t, subtest.expectedIsFound, isFound)
			require.EqualValues(t, subtest.expectedError, err)
		})
	}

}
