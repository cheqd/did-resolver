package types

import (
	"encoding/json"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
)

type DereferencingOption ResolutionOption

type DereferencingMetadata ResolutionMetadata

type DidDereferencing struct {
	ContentStream         json.RawMessage  `json:"contentStream,omitempty"`
	Metadata              cheqd.Metadata        `json:"contentMetadata,omitempty"`
	DereferencingMetadata DereferencingMetadata `json:"dereferencingMetadata,omitempty"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError ErrorType) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}
