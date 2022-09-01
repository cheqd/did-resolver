package types

import (
	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
)

type ResolutionDidDocMetadata struct {
	Created     string                 `json:"created,omitempty"`
	Updated     string                 `json:"updated,omitempty"`
	Deactivated bool                   `json:"deactivated,omitempty"`
	VersionId   string                 `json:"versionId,omitempty"`
	Resources   []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

type ResourcePreview struct {
	ResourceURI       string `json:"resourceURI"`
	CollectionId      string `json:"resourceCollectionId"`
	ResourceId        string `json:"resourceId"`
	Name              string `json:"resourceName"`
	ResourceType      string `json:"resourceType"`
	MediaType         string `json:"mediaType"`
	Created           string `json:"created"`
	Checksum          string `json:"checksum"`
	PreviousVersionId string `json:"previousVersionId"`
	NextVersionId     string `json:"nextVersionId"`
}

func NewResolutionDidDocMetadata(did string, metadata cheqd.Metadata, resources []*resource.ResourceHeader) *ResolutionDidDocMetadata {
	newMetadata := ResolutionDidDocMetadata{
		Created:     metadata.Created,
		Updated:     metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId:   metadata.VersionId,
	}
	if metadata.Resources == nil || len(resources) == 0 {
		return &newMetadata
	}
	newMetadata.Resources = NewDereferencedResourceList(did, resources).Resources
	return &newMetadata
}

func TransformToFragmentMetadata(metadata ResolutionDidDocMetadata) ResolutionDidDocMetadata {
	metadata.Resources = nil
	return metadata
}
