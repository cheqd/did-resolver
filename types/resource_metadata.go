package types

import (
	"encoding/json"
)

type ResolutionResourceMetadata struct {
	ContentMetadata *DereferencedResource     `json:"metadata,omitempty"`
	Resources       *DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
}

func (e *ResolutionResourceMetadata) MarshalJSON() ([]byte, error) {
	// If Metadata is present, use custom marshaller
	if e.ContentMetadata != nil {
		return json.Marshal(e.ContentMetadata)
	}

	// Otherwise, marshal Resources normally
	return json.Marshal(struct {
		Resources *DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
	}{
		Resources: e.Resources,
	})
}
