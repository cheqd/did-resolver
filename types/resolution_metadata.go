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
	ContentType     ContentType   `json:"contentType,omitempty"`
	ResolutionError ErrorType     `json:"error,omitempty"`
	Retrieved       string        `json:"retrieved,omitempty"`
	DidProperties   DidProperties `json:"did,omitempty"`
}

type DidProperties struct {
	DidString        string `json:"didString,omitempty"`
	MethodSpecificId string `json:"methodSpecificId,omitempty"`
	Method           string `json:"method,omitempty"`
}

type DidResolution struct {
	Did                cheqd.Did                `json:"didDocument,omitempty"`
	Metadata           ResolutionDidDocMetadata `json:"didDocumentMetadata,omitempty"`
	ResolutionMetadata ResolutionMetadata       `json:"didResolutionMetadata,omitempty"`
}

func NewResolutionMetadata(didUrl string, contentType ContentType, resolutionError ErrorType) ResolutionMetadata {
	did, _, _, _, err1 := cheqdUtils.TrySplitDIDUrl(didUrl)
	method, _, id, err2 := cheqdUtils.TrySplitDID(did)
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
