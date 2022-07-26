package types

import (
	"encoding/json"
)

type DereferencingOption ResolutionOption

type DereferencingMetadata ResolutionMetadata

type DidDereferencing struct {
	ContentStream         json.RawMessage          `json:"contentStream,omitempty"`
	Metadata              ResolutionDidDocMetadata `json:"contentMetadata,omitempty"`
	DereferencingMetadata DereferencingMetadata    `json:"dereferencingMetadata,omitempty"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError ErrorType) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}
