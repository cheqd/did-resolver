package types

import (
	"time"

	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
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
	did, _, _, _, err1 := didUtils.TrySplitDIDUrl(didUrl)
	method, _, id, err2 := didUtils.TrySplitDID(did)
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
