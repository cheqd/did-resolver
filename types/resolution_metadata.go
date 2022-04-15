package types

import (
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

const (
	ResolutionInvalidDID         = "invalidDid"
	ResolutionNotFound           = "notFound"
	ResolutionMethodNotSupported = "methodNotSupported"
)

const (
	ResolutionDIDJSONType   = "application/did+json"
	ResolutionDIDJSONLDType = "application/did+ld+json"
	ResolutionJSONLDType    = "application/ld+json"
)

const (
	DIDSchemaJSONLD = "https://ww.w3.org/ns/did/v1"
)

type ResolutionOption struct {
	Accept string `json:"accept,omitempty"`
}

type ResolutionMetadata struct {
	ContentType     string        `json:"contentType,omitempty"`
	ResolutionError string        `json:"error,omitempty"`
	Retrieved       string        `json:"retrieved,omitempty"`
	DidProperties   DidProperties `json:"did,omitempty"`
}

type DidProperties struct {
	DidString        string `json:"didString,omitempty"`
	MethodSpecificId string `json:"methodSpecificId,omitempty"`
	Method           string `json:"method,omitempty"`
}

type DidResolution struct {
	Did                cheqd.Did          `json:"didDocument,omitempty"`
	Metadata           cheqd.Metadata     `json:"didDocumentMetadata,omitempty"`
	ResolutionMetadata ResolutionMetadata `json:"didResolutionMetadata,omitempty"`
}

func NewResolutionMetadata(did string, contentType string, resolutionError string) ResolutionMetadata {
	method, _, id, _ := cheqdUtils.TrySplitDID(did)
	return ResolutionMetadata{
		ContentType:     contentType,
		ResolutionError: resolutionError,
		Retrieved:       time.Now().UTC().Format(time.RFC3339),
		DidProperties: DidProperties{
			DidString:        did,
			MethodSpecificId: id,
			Method:           method,
		},
	}
}
