package types

import resource "github.com/cheqd/cheqd-node/x/resource/types"

type DereferencedResource struct {
	Context           []string `json:"@context,omitempty"`
	CollectionId      string   `json:"collectionId,omitempty"`
	Id                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	ResourceType      string   `json:"resourceType,omitempty"`
	MediaType         string   `json:"mediaType,omitempty"`
	Created           string   `json:"created,omitempty"`
	Checksum          string   `json:"checksum,omitempty"`
	PreviousVersionId string   `json:"previousVersionId,omitempty"`
	NextVersionId     string   `json:"nextVersionId,omitempty"`
}

func NewDereferencedResource(resource *resource.ResourceHeader) *DereferencedResource {
	return &DereferencedResource{
		CollectionId:      resource.CollectionId,
		Id:                resource.Id,
		Name:              resource.Name,
		ResourceType:      resource.ResourceType,
		MediaType:         resource.MediaType,
		Created:           resource.Created,
		Checksum:          FixResourceChecksum(resource.Checksum),
		PreviousVersionId: resource.PreviousVersionId,
		NextVersionId:     resource.NextVersionId,
	}
}

func (e *DereferencedResource) AddContext(newProtocol string) {
	e.Context = AddElemToSet(e.Context, newProtocol)
}

func (e *DereferencedResource) RemoveContext() {
	e.Context = []string{}
}

func (e *DereferencedResource) GetBytes() []byte {
	return []byte{}
}

type DereferencedResourceList struct {
	Context   []string               `json:"@context,omitempty"`
	Resources []DereferencedResource `json:"resources,omitempty"`
}

func NewDereferencedResourceList(protoResources []*resource.ResourceHeader) *DereferencedResourceList {
	resourceList := []DereferencedResource{}
	for _, r := range protoResources {
		resourceList = append(resourceList, *NewDereferencedResource(r))
	}

	return &DereferencedResourceList{
		Resources: resourceList,
	}
}

func (e *DereferencedResourceList) AddContext(newProtocol string) {
	e.Context = AddElemToSet(e.Context, newProtocol)
}
func (e *DereferencedResourceList) RemoveContext()   { e.Context = []string{} }
func (e *DereferencedResourceList) GetBytes() []byte { return []byte{} }

type DereferencedResourceData []byte

func (e *DereferencedResourceData) AddContext(newProtocol string) {}
func (e *DereferencedResourceData) RemoveContext()                {}
func (e *DereferencedResourceData) GetBytes() []byte              { return []byte(*e) }
