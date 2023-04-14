package diddoc

import (
	"fmt"

	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

func transformKeyEd25519VerificationKey2018ToEd25519VerificationKey2020(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyMultibase, err := utils.Ed25519VerificationKey2018ToEd25519VerificationKey2020(
		verificationMethod.PublicKeyBase58,
	)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyBase58 = ""
	verificationMethod.Type = string(types.Ed25519VerificationKey2020)
	verificationMethod.PublicKeyMultibase = publicKeyMultibase

	return verificationMethod, nil
}

func transformKeyEd25519VerificationKey2018ToJSONWebKey2020(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyJwk, err := utils.Ed25519VerificationKey2018ToJSONWebKey2020(verificationMethod.PublicKeyBase58)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyBase58 = ""
	verificationMethod.Type = string(types.JsonWebKey2020)
	verificationMethod.PublicKeyJwk = publicKeyJwk

	return verificationMethod, nil
}

func transformKeyEd25519VerificationKey2020ToEd25519VerificationKey2018(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyBase58, err := utils.Ed25519VerificationKey2020ToEd25519VerificationKey2018(verificationMethod.PublicKeyMultibase)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyMultibase = ""
	verificationMethod.Type = string(types.Ed25519VerificationKey2018)
	verificationMethod.PublicKeyBase58 = publicKeyBase58

	return verificationMethod, nil
}

func transformKeyEd25519VerificationKey2020ToJSONWebKey2020(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyJwk, err := utils.Ed25519VerificationKey2020ToJSONWebKey2020(verificationMethod.PublicKeyMultibase)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyMultibase = ""
	verificationMethod.Type = string(types.JsonWebKey2020)
	verificationMethod.PublicKeyJwk = publicKeyJwk

	return verificationMethod, nil
}

func transformKeyJSONWebKey2020ToEd25519VerificationKey2018(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyBase58, err := utils.JSONWebKey2020ToEd25519VerificationKey2018(verificationMethod.PublicKeyJwk)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyJwk = nil
	verificationMethod.Type = string(types.Ed25519VerificationKey2018)
	verificationMethod.PublicKeyBase58 = publicKeyBase58

	return verificationMethod, nil
}

func transformKeyJSONWebKey2020ToEd25519VerificationKey2020(
	verificationMethod types.VerificationMethod,
) (types.VerificationMethod, error) {
	publicKeyMultibase, err := utils.JSONWebKey2020ToEd25519VerificationKey2020(verificationMethod.PublicKeyJwk)
	if err != nil {
		return verificationMethod, err
	}

	verificationMethod.PublicKeyJwk = nil
	verificationMethod.Type = string(types.Ed25519VerificationKey2020)
	verificationMethod.PublicKeyMultibase = publicKeyMultibase

	return verificationMethod, nil
}

func transformVerificationMethodKey(
	verificationMethod types.VerificationMethod, transformKeyType types.TransformKeyType,
) (types.VerificationMethod, error) {
	verificationMethodType := types.TransformKeyType(verificationMethod.Type)
	if verificationMethodType == transformKeyType {
		return verificationMethod, nil
	}

	switch verificationMethodType {
	case types.Ed25519VerificationKey2018:
		switch transformKeyType {
		case types.Ed25519VerificationKey2020:
			return transformKeyEd25519VerificationKey2018ToEd25519VerificationKey2020(verificationMethod)
		case types.JsonWebKey2020:
			return transformKeyEd25519VerificationKey2018ToJSONWebKey2020(verificationMethod)
		}

	case types.Ed25519VerificationKey2020:
		switch transformKeyType {
		case types.Ed25519VerificationKey2018:
			return transformKeyEd25519VerificationKey2020ToEd25519VerificationKey2018(verificationMethod)
		case types.JsonWebKey2020:
			return transformKeyEd25519VerificationKey2020ToJSONWebKey2020(verificationMethod)
		}

	case types.JsonWebKey2020:
		switch transformKeyType {
		case types.Ed25519VerificationKey2018:
			return transformKeyJSONWebKey2020ToEd25519VerificationKey2018(verificationMethod)

		case types.Ed25519VerificationKey2020:
			return transformKeyJSONWebKey2020ToEd25519VerificationKey2020(verificationMethod)
		}
	}

	return verificationMethod, fmt.Errorf("not supported transform key type")
}
