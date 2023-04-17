package types

import (
	"sort"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
)

type DereferencedDidVersionsList struct {
	Versions DidDocMetadataList `json:"versions,omitempty"`
}

func NewDereferencedDidVersionsList(versions []*didTypes.Metadata, resources []*resourceTypes.Metadata) *DereferencedDidVersionsList {
	didVersionList := DidDocMetadataList{}
	for _, version := range versions {
		didVersionList = append(didVersionList, NewResolutionDidDocMetadata("", version, resources))
	}

	// Sort by updated date or created in reverse order
	sort.Sort(didVersionList)

	return &DereferencedDidVersionsList{
		Versions: didVersionList,
	}
}

func (e *DereferencedDidVersionsList) AddContext(newProtocol string) {}
func (e *DereferencedDidVersionsList) RemoveContext()                {}
func (e *DereferencedDidVersionsList) GetBytes() []byte              { return []byte{} }
