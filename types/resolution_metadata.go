package types

import (
	"errors"
	"regexp"
	"time"
)

type ResolutionMetadata struct {
	ContentType     ContentType   `json:"contentType,omitempty" example:"application/did+ld+json"`
	ResolutionError string        `json:"error,omitempty"`
	Retrieved       string        `json:"retrieved,omitempty" example:"2021-09-01T12:00:00Z"`
	DidProperties   DidProperties `json:"did,omitempty"`
}

type DidProperties struct {
	DidString        string `json:"didString,omitempty"`
	MethodSpecificId string `json:"methodSpecificId,omitempty"`
	Method           string `json:"method,omitempty"`
}

type DidResolution struct {
	Context            string                   `json:"@context,omitempty"`
	ResolutionMetadata ResolutionMetadata       `json:"didResolutionMetadata"`
	Did                *DidDoc                  `json:"didDocument"`
	Metadata           ResolutionDidDocMetadata `json:"didDocumentMetadata"`
}

func NewResolutionMetadata(didUrl string, contentType ContentType, resolutionError string) ResolutionMetadata {
	did, _, _, _, err1 := TrySplitDIDUrl(didUrl)
	method, _, id, err2 := TrySplitDID(did)
	var didProperties DidProperties
	if err1 == nil && err2 == nil {
		didProperties = DidProperties{
			DidString:        did,
			MethodSpecificId: id,
			Method:           method,
		}
	}
	return ResolutionMetadata{
		ContentType:     contentType,
		ResolutionError: resolutionError,
		Retrieved:       time.Now().UTC().Format(time.RFC3339),
		DidProperties:   didProperties,
	}
}

func (r DidResolution) GetContentType() string {
	return string(r.ResolutionMetadata.ContentType)
}

func (r DidResolution) GetBytes() []byte {
	return []byte{}
}

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
