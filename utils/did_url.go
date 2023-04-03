package utils

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
)


// SplitDIDURLRegexp ...
// That for groups:
// Example: did:cheqd:testnet:fafdsffq11213343/path-to-s/ome-external-resource?query#key1???
// 1 - [^/?#]* - all the symbols except / and ? and # . This is the DID part                      (did:cheqd:testnet:fafdsffq11213343)
// 2 - [^?#]*  - all the symbols except ? and #. it means te section started from /, path-abempty (/path-to-s/ome-external-resource)
// 3 - \?([^#]*) - group for `query` part but with ? symbol 									  (?query)
// 4 - [^#]*     - group inside query string, match only exact query                              (query)
// 5 - #([^#]+[\$]?) - group for fragment, starts with #, includes #                              (#key1???)
// 6 - [^#]+[\$]?    - fragment only															  (key1???)
// Number of queries is not limited.
var SplitDIDURLRegexp = regexp.MustCompile(`([^/?#]*)?([^?#]*)(\?([^#]*))?(#([^#]+$))?$`)

var (
	DIDPathAbemptyRegexp = regexp.MustCompile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@]*|(%[0-9A-Fa-f]{2})*)*$`)
	DIDQueryRegexp       = regexp.MustCompile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*|(%[0-9A-Fa-f]{2})*)*$`)
	DIDFragmentRegexp    = regexp.MustCompile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*|(%[0-9A-Fa-f]{2})*)*$`)
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
	did, path, query, fragment, err := TrySplitDIDUrl(didURL)
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
	did, path, query, fragment, err := TrySplitDIDUrl(didURL)
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
	if !DIDFragmentRegexp.MatchString(fragment) {
		return fmt.Errorf("did url fragment must match the following regexp: %s", DIDFragmentRegexp)
	}
	return nil
}

func ValidateQuery(query string) error {
	if !DIDQueryRegexp.MatchString(query) {
		return fmt.Errorf("did url query must match the following regexp: %s", DIDQueryRegexp)
	}
	return nil
}

func ValidatePath(path string) error {
	if !DIDPathAbemptyRegexp.MatchString(path) {
		return fmt.Errorf("did url path abempty must match the following regexp: %s", DIDPathAbemptyRegexp)
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

func GetQuery(queryRaw string) (query string) {
	if queryRaw != "" {
		query += fmt.Sprintf("?%s", queryRaw)
	}

	return query
}

func GetFragment(fragmentId string) (fragment string) {
	if fragmentId != "" {
		fragment += url.PathEscape(fmt.Sprintf("#%s", fragmentId))
	}

	return fragment
}

// TrySplitDIDUrl Validates generic format of DIDUrl. It doesn't validate path, query and fragment content.
// Call ValidateDIDUrl for further validation.
func TrySplitDIDUrl(didURL string) (did string, path string, query string, fragment string, err error) {
	matches := SplitDIDURLRegexp.FindAllStringSubmatch(didURL, -1)

	if len(matches) != 1 {
		return "", "", "", "", errors.New("unable to split did url into did, path, query and fragment")
	}

	match := matches[0]

	return match[1], match[2], match[4], match[6], nil
}

var (
	SplitDIDRegexp     = regexp.MustCompile(`^did:([^:]+?)(:([^:]+?))?:([^:]+)$`)
	DidNamespaceRegexp = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

// TrySplitDID Validates generic format of DID. It doesn't validate method, name and id content.
// Call ValidateDID for further validation.
func TrySplitDID(did string) (method string, namespace string, id string, err error) {
	// Example: did:cheqd:testnet:base58str1ng1111
	// match [0] - the whole string
	// match [1] - cheqd                - method
	// match [2] - :testnet
	// match [3] - testnet              - namespace
	// match [4] - base58str1ng1111     - id
	matches := SplitDIDRegexp.FindAllStringSubmatch(did, -1)
	if len(matches) != 1 {
		return "", "", "", errors.New("unable to split did into method, namespace and id")
	}

	match := matches[0]
	return match[1], match[3], match[4], nil
}
