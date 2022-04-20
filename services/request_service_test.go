package services

import (
	"fmt"
	"testing"

	"github.com/cheqd/cheqd-did-resolver/types"
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
)

type MockLedgerService struct {
	Did      cheqd.Did
	Metadata cheqd.Metadata
}

func NewMockLedgerService(did cheqd.Did, metadata cheqd.Metadata) MockLedgerService {
	return MockLedgerService{
		Did:      did,
		Metadata: metadata,
	}
}

func (ls MockLedgerService) QueryDIDDoc(string) (cheqd.Did, cheqd.Metadata, bool, error) {
	isFound := true
	if ls.Did.Id == "" {
		isFound = false
	}
	return ls.Did, ls.Metadata, isFound, nil
}
func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}

func TestResolve(t *testing.T) {
	viper.SetConfigFile("../config.yaml")
	viper.ReadInConfig()
	validIdentifier := "N22KY2Dyvmuu2Pyy"
	validMethod := viper.GetString("method")
	validDIDDoc := cheqd.Did{
		Id: "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
	}
	validMetadata := cheqd.Metadata{VersionId: "test_version_id", Deactivated: false}

	subtests := []struct {
		name             string
		ledgerService    MockLedgerService
		resolutionType   string
		identifier       string
		method           string
		expectedDID      cheqd.Did
		expectedMetadata cheqd.Metadata
		expectedError    types.ResolutionError
	}{
		{
			name:             "successful resolution",
			ledgerService:    NewMockLedgerService(validDIDDoc, validMetadata),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			expectedDID:      validDIDDoc,
			expectedMetadata: validMetadata,
			expectedError:    "",
		},
		{
			name:             "DID not found",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           validMethod,
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionNotFound,
		},
		{
			name:             "invalid DID",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           validMethod,
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionInvalidDID,
		},
		{
			name:             "invalid method",
			ledgerService:    NewMockLedgerService(cheqd.Did{}, cheqd.Metadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       validIdentifier,
			method:           "not_supported_method",
			expectedDID:      cheqd.Did{},
			expectedMetadata: cheqd.Metadata{},
			expectedError:    types.ResolutionMethodNotSupported,
		},
	}

	for _, subtest := range subtests {
		t.Run(subtest.name, func(t *testing.T) {
			requestService := NewRequestService(subtest.ledgerService)
			id := "did:" + subtest.method + ":testnet:" + subtest.identifier
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
			fmt.Println(resolutionResult)
			require.EqualValues(t, subtest.expectedDID, resolutionResult.Did)
			require.EqualValues(t, subtest.expectedMetadata, resolutionResult.Metadata)
			require.EqualValues(t, subtest.resolutionType, resolutionResult.ResolutionMetadata.ContentType)
			require.EqualValues(t, subtest.expectedError, resolutionResult.ResolutionMetadata.ResolutionError)
			require.EqualValues(t, expectedDIDProperties, resolutionResult.ResolutionMetadata.DidProperties)
			require.Empty(t, err)
		})
	}

}