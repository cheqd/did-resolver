package types

import (
	"time"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
)

type ResolutionResourceMetadata struct {
	Created           *time.Time               `json:"created,omitempty" example:"2021-09-01T12:00:00Z"`
	Updated           *time.Time               `json:"updated,omitempty" example:"2021-09-10T12:00:00Z"`
	Deactivated       bool                     `json:"deactivated,omitempty" example:"false"`
	VersionId         string                   `json:"versionId,omitempty" example:"284f297b-b6e3-4ffa-9172-bc3bb904e286"`
	NextVersionId     string                   `json:"nextVersionId,omitempty" example:"3f3111af-dfe6-411f-adc9-02af59716ddb"`
	PreviousVersionId string                   `json:"previousVersionId,omitempty" example:"139445af-4281-4453-b05a-ec9a8931c1f9"`
	Resources         DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
}

func NewResolutionResourceMetadata(did string, metadata *didTypes.Metadata, resources []*resourceTypes.Metadata) ResolutionResourceMetadata {
	created := toTime(metadata.Created)
	updated := toTime(metadata.Updated)

	newMetadata := ResolutionResourceMetadata{
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

func (e *ResolutionResourceMetadata) AddContext(newProtocol string) {}
func (e *ResolutionResourceMetadata) RemoveContext()                {}
func (e *ResolutionResourceMetadata) GetBytes() []byte              { return []byte{} }
func (e *ResolutionResourceMetadata) GetContentType() string        { return "" }
func (e *ResolutionResourceMetadata) IsRedirect() bool              { return false }
