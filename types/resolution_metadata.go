package types

import (
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

type ResolutionOption struct {
	Accept ContentType `json:"accept,omitempty"`
}

type ResolutionMetadata struct {
	ContentType     ContentType     `json:"contentType,omitempty"`
	ResolutionError ResolutionError `json:"error,omitempty"`
	Retrieved       string          `json:"retrieved,omitempty"`
	DidProperties   DidProperties   `json:"did,omitempty"`
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

func NewResolutionMetadata(did string, contentType ContentType, resolutionError ResolutionError) ResolutionMetadata {
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
