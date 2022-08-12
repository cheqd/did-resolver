package types

import resource "github.com/cheqd/cheqd-node/x/resource/types"

type DereferencedResource struct {
	Context           []string `json:"@context,omitempty"`
	CollectionId      string   `json:"collectionId,omitempty"`
	Id                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	ResourceType      string   `json:"resourceType,omitempty"`
	MediaType         string   `json:"mediaType,omitempty"`
	Created           string   `json:"created,omitempty"`
	Checksum          string   `json:"checksum,omitempty"`
	PreviousVersionId string   `json:"previousVersionId,omitempty"`
	NextVersionId     string   `json:"nextVersionId,omitempty"`
}

func NewDereferencedResource(context []string, resource *resource.ResourceHeader) DereferencedResource {
	return DereferencedResource{
		Context:           context,
		CollectionId:      resource.CollectionId,
		Id:                resource.Id,
		Name:              resource.Name,
		ResourceType:      resource.ResourceType,
		MediaType:         resource.MediaType,
		Created:           resource.Created,
		Checksum:          FixResourceChecksum(resource.Checksum),
		PreviousVersionId: resource.PreviousVersionId,
		NextVersionId:     resource.NextVersionId,
	}
}
