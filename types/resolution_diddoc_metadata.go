package types

import (
	"time"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ResolutionDidDocMetadata struct {
	Created           *time.Time               `json:"created,omitempty" example:"2021-09-01T12:00:00Z"`
	Updated           *time.Time               `json:"updated,omitempty" example:"2021-09-10T12:00:00Z"`
	Deactivated       bool                     `json:"deactivated,omitempty" example:"false"`
	VersionId         string                   `json:"versionId,omitempty" example:"284f297b-b6e3-4ffa-9172-bc3bb904e286"`
	NextVersionId     string                   `json:"nextVersionId,omitempty" example:"3f3111af-dfe6-411f-adc9-02af59716ddb"`
	PreviousVersionId string                   `json:"previousVersionId,omitempty" example:"139445af-4281-4453-b05a-ec9a8931c1f9"`
	Resources         DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionDidDocMetadata(did string, metadata *didTypes.Metadata, resources []*resourceTypes.Metadata) ResolutionDidDocMetadata {
	created := toTime(metadata.Created)
	updated := toTime(metadata.Updated)

	newMetadata := ResolutionDidDocMetadata{
		Created:           created,
		Updated:           updated,
		Deactivated:       metadata.Deactivated,
		VersionId:         metadata.VersionId,
		NextVersionId:     metadata.NextVersionId,
		PreviousVersionId: metadata.PreviousVersionId,
	}

	if len(resources) == 0 {
		return newMetadata
	}

	newMetadata.Resources = NewDereferencedResourceListStruct(did, resources).Resources
	return newMetadata
}

func TransformToFragmentMetadata(metadata ResolutionDidDocMetadata) ResolutionDidDocMetadata {
	metadata.Resources = nil
	return metadata
}

func (e *ResolutionDidDocMetadata) AddContext(newProtocol string) {}
func (e *ResolutionDidDocMetadata) RemoveContext()                {}
func (e *ResolutionDidDocMetadata) GetBytes() []byte              { return []byte{} }
func (e *ResolutionDidDocMetadata) GetContentType() string        { return "" }
func (e *ResolutionDidDocMetadata) IsRedirect() bool              { return false }

func toTime(value *timestamppb.Timestamp) (result *time.Time) {
	if value == nil || value.AsTime().IsZero() {
		result = nil
	} else {
		value, _ := time.Parse(time.RFC3339, value.AsTime().Format(time.RFC3339))
		result = &value
	}

	return result
}
