package types

import (
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
)

type DereferencingOption ResolutionOption

type DereferencingMetadata ResolutionMetadata

type DidDereferencing struct {
	ContentStream      interface{}           `json:"contentStream,omitempty"`
	Metadata           cheqd.Metadata        `json:"contentMetadata,omitempty"`
	ResolutionMetadata DereferencingMetadata `json:"DereferencingMetadata,omitempty"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError ResolutionError) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}
