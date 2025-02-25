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

func (e *ResolutionResourceMetadata) UnmarshalJSON(data []byte) error {
	// Define a temporary structure to assist with unmarshalling
	var aux struct {
		Resources *DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
	}

	// First, try to unmarshal into ContentMetadata
	if err := json.Unmarshal(data, &e.ContentMetadata); err == nil && e.ContentMetadata.CollectionId != "" {
		return nil // Successfully unmarshalled into ContentMetadata, return early
	}

	// If ContentMetadata is nil, try to unmarshal into Resources
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Assign the extracted Resources
	e.Resources = aux.Resources
	e.ContentMetadata = nil
	return nil
}
