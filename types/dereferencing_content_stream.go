package types

type DereferencedResource struct {
	Context           []string `json:"@context,omitempty"`
	CollectionId      string   `json:"collectionId,omitempty"`
	Id                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	ResourceType      string   `json:"resourceType,omitempty"`
	MediaType         string   `json:"mediaType,omitempty"`
	Created           string   `json:"created,omitempty"`
	Checksum          []byte   `json:"checksum,omitempty"`
	PreviousVersionId string   `json:"previousVersionId,omitempty"`
	NextVersionId     string   `json:"nextVersionId,omitempty"`
	Data              []byte   `json:"data,omitempty"`
}
