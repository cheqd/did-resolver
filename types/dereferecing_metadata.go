package types

import (
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	"google.golang.org/protobuf/runtime/protoiface"
)

type DereferencingOption ResolutionOption

type DereferencingMetadata ResolutionMetadata

type ContentStream protoiface.MessageV1

type DidDereferencing struct {
	ContentStream         protoiface.MessageV1  `json:"contentStream,omitempty"`
	Metadata              cheqd.Metadata        `json:"contentMetadata,omitempty"`
	DereferencingMetadata DereferencingMetadata `json:"dereferencingMetadata,omitempty"`
}

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError ErrorType) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}
