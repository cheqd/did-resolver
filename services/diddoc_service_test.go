package services

import (
	"testing"

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
