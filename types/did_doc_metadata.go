package types

import (
	did "github.com/cheqd/cheqd-node/api/cheqd/did/v2"
	resource "github.com/cheqd/cheqd-node/api/cheqd/resource/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	newMetadata := ResolutionDidDocMetadata{
		Created:     created,
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
