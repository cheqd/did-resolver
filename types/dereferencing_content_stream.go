package types

import resource "github.com/cheqd/cheqd-node/x/resource/types"

type DereferencedResource struct {
	ResourceURI       string  `json:"resourceURI"`
	CollectionId      string  `json:"resourceCollectionId"`
	ResourceId        string  `json:"resourceId"`
	Name              string  `json:"resourceName"`
	ResourceType      string  `json:"resourceType"`
	MediaType         string  `json:"mediaType"`
	Created           string  `json:"created"`
	Checksum          string  `json:"checksum"`
	PreviousVersionId *string `json:"previousVersionId"`
	NextVersionId     *string `json:"nextVersionId"`
}

func NewDereferencedResource(did string, resource *resource.ResourceHeader) *DereferencedResource {
	var previousVersionId, nextVersionId *string
	if resource.PreviousVersionId != "" {
		previousVersionId = &resource.PreviousVersionId
	}
	if resource.NextVersionId != "" {
		nextVersionId = &resource.NextVersionId
	}
	return &DereferencedResource{
		ResourceURI:       did + RESOURCE_PATH + resource.Id,
		CollectionId:      resource.CollectionId,
		ResourceId:        resource.Id,
		Name:              resource.Name,
		ResourceType:      resource.ResourceType,
		MediaType:         resource.MediaType,
		Created:           resource.Created,
		Checksum:          FixResourceChecksum(resource.Checksum),
		PreviousVersionId: previousVersionId,
		NextVersionId:     nextVersionId,
	}
}

type DereferencedResourceList struct {
	Resources []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewDereferencedResourceList(did string, protoResources []*resource.ResourceHeader) *DereferencedResourceList {
	resourceList := []DereferencedResource{}
	for _, r := range protoResources {
		resourceList = append(resourceList, *NewDereferencedResource(did, r))
	}

	return &DereferencedResourceList{
		Resources: resourceList,
	}
}

func (e *DereferencedResourceList) AddContext(newProtocol string) {}
func (e *DereferencedResourceList) RemoveContext()                {}
func (e *DereferencedResourceList) GetBytes() []byte              { return []byte{} }

type DereferencedResourceData []byte

func (e *DereferencedResourceData) AddContext(newProtocol string) {}
func (e *DereferencedResourceData) RemoveContext()                {}
func (e *DereferencedResourceData) GetBytes() []byte              { return []byte(*e) }
