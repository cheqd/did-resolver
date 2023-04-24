package utils

import (
	"encoding/base64"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/multiformats/go-multibase"
)

func Ed25519VerificationKey2018ToEd25519VerificationKey2020(publicKeyBase58 string) (string, error) {
	// get public key from Ed25519VerificationKey2018 key
	pubKey, err := base58.Decode(publicKeyBase58)
	if err != nil {
		return "", err
	}

	// generate Ed25519VerificationKey2020 key
	return GenerateEd25519VerificationKey2020(pubKey)
}

func Ed25519VerificationKey2018ToJSONWebKey2020(publicKeyBase58 string) (interface{}, error) {
	// get public key from Ed25519VerificationKey2018 key
	pubKey, err := base58.Decode(publicKeyBase58)
	if err != nil {
		return "", err
	}

	// generate GenerateJSONWebKey2020 key
	return GenerateJSONWebKey2020(pubKey)
}

func Ed25519VerificationKey2020ToEd25519VerificationKey2018(publicKeyMultibase string) (string, error) {
	// get public key from Ed25519VerificationKey2020 key
	encoding, pubKey, err := multibase.Decode(publicKeyMultibase)
	if encoding != multibase.Base58BTC {
		return "", fmt.Errorf("Only Base58BTC encoding is supported")
	}
	if err != nil {
		return "", err
	}

	// generate Ed25519VerificationKey2018 key
	return GenerateEd25519VerificationKey2018(pubKey), nil
}

func Ed25519VerificationKey2020ToJSONWebKey2020(publicKeyMultibase string) (interface{}, error) {
	// get public key from Ed25519VerificationKey2020 key
	encoding, pubKey, err := multibase.Decode(publicKeyMultibase)
	if encoding != multibase.Base58BTC {
		return "", fmt.Errorf("Only Base58BTC encoding is supported")
	}
	if err != nil {
		return "", err
	}

	// generate JSONWebKey2020 key
	return GenerateJSONWebKey2020(pubKey)
}

func JSONWebKey2020ToEd25519VerificationKey2018(publicKeyJwk interface{}) (string, error) {
	// get the public key from JSONWebKey2020
	jwk := publicKeyJwk.(map[string]interface{})
	pubKey, err := base64.RawURLEncoding.DecodeString(jwk["x"].(string))
	if err != nil {
		return "", err
	}

	// generate Ed25519VerificationKey2018 key
	return GenerateEd25519VerificationKey2018(pubKey), nil
}

func JSONWebKey2020ToEd25519VerificationKey2020(publicKeyJwk interface{}) (string, error) {
	// get the public key from JSONWebKey2020
	jwk := publicKeyJwk.(map[string]interface{})
	pubKey, err := base64.RawURLEncoding.DecodeString(jwk["x"].(string))
	if err != nil {
		return "", err
	}

	// generate Ed25519VerificationKey2020 key
	return GenerateEd25519VerificationKey2020(pubKey)
}
