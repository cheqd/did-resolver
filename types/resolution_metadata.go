package types

import (
	"time"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
)

const (
	ResolutionInvalidDID                 = "invalidDid"
	ResolutionNotFound                   = "notFound"
	ResolutionRepresentationNotSupported = "representationNotSupported"
)

const (
	ResolutionJSONType   = "application/json"
	ResolutionJSONLDType = "application/ld+json"
)

type ResolutionOption struct {
	Accept string `json:"accept,omitempty"`
}

type ResolutionMetadata struct {
	ContentType     string `json:"contentType,omitempty"`
	ResolutionError string `json:"error,omitempty"`
	Retrieved       string `json:"retrieved,omitempty"`
}

type DidResolution struct {
	Did                cheqd.Did          `json:"didDocument,omitempty"`
	Metadata           cheqd.Metadata     `json:"didDocumentMetadata,omitempty"`
	ResolutionMetadata ResolutionMetadata `json:"didResolutionMetadata,omitempty"`
}

func NewResolutionMetadata(contentType string, resolutionError string) ResolutionMetadata {
	return ResolutionMetadata{
		contentType,
		resolutionError,
		time.Now().UTC().Format(time.RFC3339),
	}
}
