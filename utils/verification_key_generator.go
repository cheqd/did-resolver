package utils

import (
	"crypto/ed25519"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
)

func GenerateEd25519VerificationKey2018(publicKey ed25519.PublicKey) string {
	return base58.Encode(publicKey)
}

func GenerateEd25519VerificationKey2020(publicKey ed25519.PublicKey) (string, error) {
	publicKeyMultibaseBytes := []byte{0xed, 0x01}
	publicKeyMultibaseBytes = append(publicKeyMultibaseBytes, publicKey...)

	return multibase.Encode(multibase.Base58BTC, publicKeyMultibaseBytes)
}

func GenerateJSONWebKey2020(publicKey ed25519.PublicKey) (jwk.Key, error) {
	return jwk.New(publicKey)
}
