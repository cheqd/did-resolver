package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cheqd/did-resolver/types"
)

func MustSplitDID(did string) (method string, namespace string, id string) {
	method, namespace, id, err := types.TrySplitDID(did)
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

func ReplaceDidInDidURL(didURL string, oldDid string, newDid string) string {
	did, path, query, fragment := MustSplitDIDUrl(didURL)
	if did == oldDid {
		did = newDid
	}

	return JoinDIDUrl(did, path, query, fragment)
}

func ReplaceDidInDidURLList(didURLList []string, oldDid string, newDid string) []string {
	res := make([]string, len(didURLList))

	for i := range didURLList {
		res[i] = ReplaceDidInDidURL(didURLList[i], oldDid, newDid)
	}

	return res
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDID(did string, method string, allowedNamespaces []string) error {
	sMethod, sNamespace, sUniqueID, err := types.TrySplitDID(did)
	if err != nil {
		return err
	}

	// check method
	if method != "" && method != sMethod {
		return fmt.Errorf("did method must be: %s", method)
	}

	// check namespaces
	if !types.DidNamespaceRegexp.MatchString(sNamespace) {
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

// Normalization

func NormalizeDID(did string) string {
	method, namespace, id := MustSplitDID(did)
	id = NormalizeID(id)
	return JoinDID(method, namespace, id)
}

func NormalizeDIDList(didList []string) []string {
	if didList == nil {
		return nil
	}
	newDIDs := []string{}
	for _, did := range didList {
		newDIDs = append(newDIDs, NormalizeDID(did))
	}
	return newDIDs
}
