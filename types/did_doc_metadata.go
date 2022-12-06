package types

import (
	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
)

type ResolutionDidDocMetadata struct {
	Created     string                 `json:"created,omitempty" example:"2021-09-01T12:00:00Z"`
	Updated     string                 `json:"updated,omitempty" example:"2021-09-10T12:00:00Z"`
	Deactivated bool                   `json:"deactivated,omitempty" example:"false"`
	VersionId   string                 `json:"versionId,omitempty" example:"4979BAF49599FEF0BAD5ED0849FDD708156761EBBC8EBE78D0907F8BECC9CB2E"`
	Resources   []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionDidDocMetadata(did string, metadata did.Metadata, resources []*resource.Metadata) ResolutionDidDocMetadata {
	newMetadata := ResolutionDidDocMetadata{
		Created:     metadata.Created,
		Updated:     metadata.Updated,
		Deactivated: metadata.Deactivated,
		VersionId:   metadata.VersionId,
	}

	if len(resources) == 0 {
		return newMetadata
	}

	newMetadata.Resources = NewDereferencedResourceList(did, resources).Resources
	return newMetadata
}

func TransformToFragmentMetadata(metadata ResolutionDidDocMetadata) ResolutionDidDocMetadata {
	metadata.Resources = nil
	return metadata
}
