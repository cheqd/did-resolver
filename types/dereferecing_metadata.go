package types

import (
	"sort"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	"github.com/cheqd/did-resolver/utils"
)

type (
	DereferencingMetadata ResolutionMetadata
	DidDereferencing      struct {
		Context               string                   `json:"@context,omitempty" example:"https://w3id.org/did-resolution/v1"`
		DereferencingMetadata DereferencingMetadata    `json:"dereferencingMetadata"`
		ContentStream         ContentStreamI           `json:"contentStream"`
		Metadata              ResolutionDidDocMetadata `json:"contentMetadata"`
	}
)

func NewDereferencingMetadata(did string, contentType ContentType, resolutionError string) DereferencingMetadata {
	return DereferencingMetadata(NewResolutionMetadata(did, contentType, resolutionError))
}

// Interface implementation

func (d DidDereferencing) GetContentType() string {
	return string(d.DereferencingMetadata.ContentType)
}

func (d DidDereferencing) GetBytes() []byte {
	if d.ContentStream == nil {
		return []byte{}
	}
	return d.ContentStream.GetBytes()
}

func (r DidDereferencing) IsRedirect() bool {
	return false
}

// end of Interface implementation

type ResourceDereferencing struct {
	Context               string                     `json:"@context,omitempty" example:"https://w3id.org/did-resolution/v1"`
	DereferencingMetadata DereferencingMetadata      `json:"dereferencingMetadata"`
	ContentStream         ContentStreamI             `json:"contentStream"`
	Metadata              ResolutionResourceMetadata `json:"contentMetadata"`
}

// Interface implementation

func (d ResourceDereferencing) GetContentType() string {
	return string(d.DereferencingMetadata.ContentType)
}

func (d ResourceDereferencing) GetBytes() []byte {
	if d.ContentStream == nil {
		return []byte{}
	}
	return d.ContentStream.GetBytes()
}

func (r ResourceDereferencing) IsRedirect() bool {
	return false
}

type DidDocMetadataList []ResolutionDidDocMetadata

func (dd DidDocMetadataList) Len() int {
	return len(dd)
}

// Sort in reverse order
func (dd DidDocMetadataList) Less(i, j int) bool {
	if dd[i].Updated == nil {
		return false
	}
	if dd[j].Updated == nil {
		return true
	}
	return dd[i].Updated.After(*dd[j].Updated)
}

func (dd DidDocMetadataList) Swap(i, j int) {
	dd[i], dd[j] = dd[j], dd[i]
}

type DereferencedDidVersionsList struct {
	Versions DidDocMetadataList `json:"versions,omitempty"`
}

func NewDereferencedDidVersionsList(versions []*didTypes.Metadata) *DereferencedDidVersionsList {
	didVersionList := DidDocMetadataList{}
	for _, version := range versions {
		didVersionList = append(didVersionList, NewResolutionDidDocMetadata("", version, nil))
	}

	return &DereferencedDidVersionsList{
		Versions: didVersionList,
	}
}

func (e *DereferencedDidVersionsList) AddContext(newProtocol string) {}
func (e *DereferencedDidVersionsList) RemoveContext()                {}
func (e *DereferencedDidVersionsList) GetBytes() []byte              { return []byte{} }

// Returns VersionId if there is a version before the given time
// Otherwise NotFound error
func (e DereferencedDidVersionsList) FindBeforeTime(stime string) (string, error) {
	search_time, err := utils.ParseFromStringTimeToGoTime(stime)
	if err != nil {
		return "", err
	}
	// Firstly - sort versions by Updated time
	versions := e.Versions
	sort.Sort(versions)
	for _, version := range versions {
		if version.Updated != nil && version.Updated.Before(search_time) {
			return version.VersionId, nil
		}
		if version.Updated == nil && version.Created.Before(search_time) {
			return version.VersionId, nil
		}
	}
	return "", nil
}
