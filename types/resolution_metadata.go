package types

import (
	"errors"
	"time"

	"github.com/cheqd/did-resolver/utils"
)

type ResolutionMetadata struct {
	ContentType     ContentType   `json:"contentType,omitempty" example:"application/ld+json"`
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
	Context            string                    `json:"@context,omitempty"`
	ResolutionMetadata ResolutionMetadata        `json:"didResolutionMetadata"`
	Did                *DidDoc                   `json:"didDocument,omitempty"`
	Metadata           *ResolutionDidDocMetadata `json:"didDocumentMetadata,omitempty"`
}

// Interface implementation

func (r DidResolution) GetContentType() string {
	return string(r.ResolutionMetadata.ContentType)
}

func (r DidResolution) GetBytes() []byte {
	return []byte{}
}

func (r DidResolution) IsRedirect() bool {
	return false
}

func (r DidResolution) GetServiceByName(serviceName string) (string, error) {
	if r.Did == nil {
		return "", errors.New("did document is nil")
	}
	return r.Did.GetServiceByName(serviceName)
}

// end of Interface implementation

func NewResolutionMetadata(didUrl string, contentType ContentType, resolutionError string) ResolutionMetadata {
	did, _, _, _, err1 := utils.TrySplitDIDUrl(didUrl)
	method, _, id, err2 := utils.TrySplitDID(did)
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
