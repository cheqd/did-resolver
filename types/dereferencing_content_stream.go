package types

import (
	"time"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
)

type DereferencedResource struct {
	ResourceURI       string     `json:"resourceURI" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47/resources/398cee0a-efac-4643-9f4c-74c48c72a14b"`
	CollectionId      string     `json:"resourceCollectionId" example:"55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"`
	ResourceId        string     `json:"resourceId" example:"398cee0a-efac-4643-9f4c-74c48c72a14b"`
	Name              string     `json:"resourceName" example:"Image Resource"`
	ResourceType      string     `json:"resourceType" example:"Image"`
	MediaType         string     `json:"mediaType" example:"image/png"`
	Created           *time.Time `json:"created" example:"2021-09-01T12:00:00Z"`
	Checksum          string     `json:"checksum" example:"a95380f460e63ad939541a57aecbfd795fcd37c6d78ee86c885340e33a91b559"`
	PreviousVersionId *string    `json:"previousVersionId" example:"ad7a8442-3531-46eb-a024-53953ec6e4ff"`
	NextVersionId     *string    `json:"nextVersionId" example:"d4829ac7-4566-478c-a408-b44767eddadc"`
}

func NewDereferencedResource(did string, resource *resourceTypes.Metadata) *DereferencedResource {
	var previousVersionId, nextVersionId *string
	if resource.PreviousVersionId != "" {
		previousVersionId = &resource.PreviousVersionId
	}
	if resource.NextVersionId != "" {
		nextVersionId = &resource.NextVersionId
	}
	created := resource.Created.AsTime()
	return &DereferencedResource{
		ResourceURI:       did + RESOURCE_PATH + resource.Id,
		CollectionId:      resource.CollectionId,
		ResourceId:        resource.Id,
		Name:              resource.Name,
		ResourceType:      resource.ResourceType,
		MediaType:         resource.MediaType,
		Created:           &created,
		Checksum:          resource.Checksum,
		PreviousVersionId: previousVersionId,
		NextVersionId:     nextVersionId,
	}
}

type DereferencedResourceList struct {
	Resources []DereferencedResource `json:"linkedResourceMetadata,omitempty"`
}

func NewDereferencedResourceList(did string, protoResources []*resourceTypes.Metadata) *DereferencedResourceList {
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
