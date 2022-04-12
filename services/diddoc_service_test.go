package services

import (
	"fmt"
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	//nolint
)

func TestMarshallDID(t *testing.T) {
	didDocService := DIDDocService{}
	verificationMethod1 := cheqd.VerificationMethod{
		Id:                 "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey",
		Type:               "Ed25519VerificationKey2020",
		Controller:         "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
		PublicKeyMultibase: "zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf",
	}

	verificationMethod2 := cheqd.VerificationMethod{
		Id:         "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey",
		Type:       "JsonWebKey2020",
		Controller: "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
		PublicKeyJwk: []*cheqd.KeyValuePair{
			{"kty", "OKP"},
			{"crv", "Ed25519"},
			{"x", "VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ"},
		},
	}
	didDoc := cheqd.Did{
		Context:            []string{"test"},
		Id:                 "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
		VerificationMethod: []*cheqd.VerificationMethod{&verificationMethod1, &verificationMethod2},
	}

	expectedDID := "{\"@context\":[\"test\"],\"id\":\"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\"verificationMethod\":[{\"controller\":\"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\"id\":\"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey\",\"publicKeyMultibase\":\"zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf\",\"type\":\"Ed25519VerificationKey2020\"},{\"controller\":\"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\"id\":\"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey\",\"publicKeyJwk\":{\"crv\":\"Ed25519\",\"kty\":\"OKP\",\"x\":\"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"},\"type\":\"JsonWebKey2020\"}]}"

	jsonDID, err := didDocService.MarshallDID(didDoc)

	fmt.Println(jsonDID)
	require.EqualValues(t, jsonDID, expectedDID)
	require.Empty(t, err)
}
