package utils

import (
	"errors"
	"fmt"

	"github.com/cheqd/did-resolver/types"
	"github.com/google/uuid"
)

func IsValidResourceId(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func ValidateV1ID(id string) error {
	isValidId := len(id) == 16 && IsValidBase58(id) ||
		len(id) == 32 && IsValidBase58(id) ||
		IsValidUUID(id)

	if !isValidId {
		return errors.New("unique id should be one of: 16 symbols base58 string, 32 symbols base58 string, or UUID")
	}

	return nil
}

func IsValidV1ID(id string) bool {
	err := ValidateV1ID(id)
	return err == nil
}

func MustSplitDIDUrl(didURL string) (did string, path string, query string, fragment string) {
	did, path, query, fragment, err := types.TrySplitDIDUrl(didURL)
	if err != nil {
		panic(err.Error())
	}
	return
}

func JoinDIDUrl(did string, path string, query string, fragment string) string {
	res := did + path

	if query != "" {
		res = res + "?" + query
	}

	if fragment != "" {
		res = res + "#" + fragment
	}

	return res
}

// ValidateDIDUrl checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDIDUrl(didURL string, method string, allowedNamespaces []string) error {
	did, path, query, fragment, err := types.TrySplitDIDUrl(didURL)
	if err != nil {
		return err
	}

	// Validate DID
	err = ValidateDID(did, method, allowedNamespaces)
	if err != nil {
		return err
	}
	// Validate path
	err = ValidatePath(path)
	if err != nil {
		return err
	}
	// Validate query
	err = ValidateQuery(query)
	if err != nil {
		return err
	}
	// Validate fragment
	err = ValidateFragment(fragment)
	if err != nil {
		return err
	}

	return nil
}

func ValidateFragment(fragment string) error {
	if !types.DIDFragmentRegexp.MatchString(fragment) {
		return fmt.Errorf("did url fragment must match the following regexp: %s", types.DIDFragmentRegexp)
	}
	return nil
}

func ValidateQuery(query string) error {
	if !types.DIDQueryRegexp.MatchString(query) {
		return fmt.Errorf("did url query must match the following regexp: %s", types.DIDQueryRegexp)
	}
	return nil
}

func ValidatePath(path string) error {
	if !types.DIDPathAbemptyRegexp.MatchString(path) {
		return fmt.Errorf("did url path abempty must match the following regexp: %s", types.DIDPathAbemptyRegexp)
	}
	return nil
}

func IsValidDIDUrl(didURL string, method string, allowedNamespaces []string) bool {
	err := ValidateDIDUrl(didURL, method, allowedNamespaces)

	return nil == err
}

func IsMigrationNeeded(id string) bool {
	return IsValidV1ID(id)
}
