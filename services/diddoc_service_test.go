package services

import (
	"testing"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/stretchr/testify/require"
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
			{Key: "kty", Value: "OKP"},
			{Key: "crv", Value: "Ed25519"},
			{Key: "x", Value: "VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ"},
		},
	}
	didDoc := cheqd.Did{
		Context:            []string{"test"},
		Id:                 "did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue",
		VerificationMethod: []*cheqd.VerificationMethod{&verificationMethod1, &verificationMethod2},
	}

	expectedDID := "{\n" +
		"  \"@context\": [\n" +
		"    \"test\"\n" +
		"  ],\n" +
		"  \"id\": \"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\n" +
		"  \"verificationMethod\": [\n" +
		"    {\n" +
		"      \"id\": \"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey\",\n" +
		"      \"type\": \"Ed25519VerificationKey2020\",\n" +
		"      \"controller\": \"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\n" +
		"      \"publicKeyMultibase\": \"zAKJP3f7BD6W4iWEQ9jwndVTCBq8ua2Utt8EEjJ6Vxsf\"\n" +
		"    },\n" +
		"    {\n" +
		"      \"id\": \"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue#verkey\",\n" +
		"      \"type\": \"JsonWebKey2020\",\n" +
		"      \"controller\": \"did:cheqd:mainnet:N22KY2Dyvmuu2PyyqSFKue\",\n" +
		"      \"publicKeyJwk\": {\n" +
		"        \"crv\": \"Ed25519\",\n" +
		"        \"kty\": \"OKP\",\n" +
		"        \"x\": \"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"\n" +
		"      }\n" +
		"    }\n" +
		"  ]\n" +
		"}"
	jsonDID, err := didDocService.MarshallDID(didDoc)

	require.EqualValues(t, expectedDID, jsonDID)
	require.Empty(t, err)
}
