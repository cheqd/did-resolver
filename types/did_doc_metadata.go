package types

import (
	"time"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResolutionDidDocMetadata struct {
	Created     *time.Time             `json:"created,omitempty" example:"2021-09-01T12:00:00Z"`
	Updated     *time.Time             `json:"updated,omitempty" example:"2021-09-10T12:00:00Z"`
	Deactivated bool                   `json:"deactivated,omitempty" example:"false"`
	VersionId   string                 `json:"versionId,omitempty" example:"4979BAF49599FEF0BAD5ED0849FDD708156761EBBC8EBE78D0907F8BECC9CB2E"`
	Resources   []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionDidDocMetadata(did string, metadata *didTypes.Metadata, resources []*resourceTypes.Metadata) ResolutionDidDocMetadata {
	created := toTime(metadata.Created)
	updated := toTime(metadata.Updated)

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

func (e *ResolutionDidDocMetadata) AddContext(newProtocol string) {}
func (e *ResolutionDidDocMetadata) RemoveContext()                {}
func (e *ResolutionDidDocMetadata) GetBytes() []byte              { return []byte{} }

func toTime(value *timestamppb.Timestamp) (result *time.Time) {
	if value.AsTime().IsZero() {
		result = nil
	} else {
		value := value.AsTime()
		result = &value
	}

	return result
}
