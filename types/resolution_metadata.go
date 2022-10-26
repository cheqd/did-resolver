package types

import (
	"time"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

type ResolutionMetadata struct {
	ContentType     ContentType   `json:"contentType,omitempty"`
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
	Context            string                   `json:"@context,omitempty"`
	ResolutionMetadata ResolutionMetadata       `json:"didResolutionMetadata"`
	Did                *DidDoc                  `json:"didDocument"`
	Metadata           ResolutionDidDocMetadata `json:"didDocumentMetadata"`
}

func NewResolutionMetadata(didUrl string, contentType ContentType, resolutionError string) ResolutionMetadata {
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

func (r DidResolution) GetContentType() string {
	return string(r.ResolutionMetadata.ContentType)
}

func (r DidResolution) GetBytes() []byte {
	return []byte{}
}
