package utils

import (
	"errors"
	"fmt"
	"strings"
)

func MustSplitDID(did string) (method string, namespace string, id string) {
	method, namespace, id, err := TrySplitDID(did)
	if err != nil {
		panic(err.Error())
	}
	return
}

func JoinDID(method, namespace, id string) string {
	res := "did:" + method

	if namespace != "" {
		res = res + ":" + namespace
	}

	return res + ":" + id
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDID(did string, method string, allowedNamespaces []string) error {
	sMethod, sNamespace, sUniqueID, err := TrySplitDID(did)
	if err != nil {
		return err
	}

	// check method
	if method != "" && method != sMethod {
		return fmt.Errorf("did method must be: %s", method)
	}

	// check namespaces
	if !DidNamespaceRegexp.MatchString(sNamespace) {
		return errors.New("invalid did namespace")
	}

	if len(allowedNamespaces) > 0 && !Contains(allowedNamespaces, sNamespace) {
		return fmt.Errorf("did namespace must be one of: %s", strings.Join(allowedNamespaces, ", "))
	}

	// check unique-id
	err = ValidateID(sUniqueID)
	if err != nil {
		return err
	}

	return err
}

func IsValidDID(did string, method string, allowedNamespaces []string) bool {
	err := ValidateDID(did, method, allowedNamespaces)
	return err == nil
}
