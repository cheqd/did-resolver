package types

import (
	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
)

type ResolutionDidDocMetadata struct {
	Created     string                 `json:"created,omitempty"`
	Updated     string                 `json:"updated,omitempty"`
	Deactivated bool                   `json:"deactivated,omitempty"`
	VersionId   string                 `json:"versionId,omitempty"`
	Resources   []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionDidDocMetadata(did string, metadata did.Metadata, resources []*resource.Metadata) ResolutionDidDocMetadata {
	newMetadata := ResolutionDidDocMetadata{
		Created:     metadata.Created,
		Updated:     metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId:   metadata.VersionId,
	}

	newMetadata.Resources = NewDereferencedResourceList(did, resources).Resources
	return newMetadata
}

func TransformToFragmentMetadata(metadata ResolutionDidDocMetadata) ResolutionDidDocMetadata {
	metadata.Resources = nil
	return metadata
}
