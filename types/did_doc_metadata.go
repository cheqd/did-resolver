package types

import (
	"time"

	did "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
)

// Changed "Created time.Time" to "Create *time.Time".
// It needs to skip Created field when is empty.
type ResolutionDidDocMetadata struct {
	Created     *time.Time             `json:"created,omitempty" example:"2021-09-01T12:00:00Z"`
	Updated     *time.Time             `json:"updated,omitempty" example:"2021-09-10T12:00:00Z"`
	Deactivated bool                   `json:"deactivated,omitempty" example:"false"`
	VersionId   string                 `json:"versionId,omitempty" example:"4979BAF49599FEF0BAD5ED0849FDD708156761EBBC8EBE78D0907F8BECC9CB2E"`
	Resources   []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionDidDocMetadata(did string, metadata did.Metadata, resources []*resource.Metadata) ResolutionDidDocMetadata {
	created := &metadata.Created
	if created.IsZero() {
		created = nil
	}

	updated := metadata.Updated
	if updated.IsZero() {
		updated = nil
	}

	newMetadata := ResolutionDidDocMetadata{
		Created:     created,
		Updated:     updated,
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
