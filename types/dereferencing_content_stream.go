package types

import (
	"sort"
	"time"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/utils"
)

type DereferencedResource struct {
	ResourceURI       string     `json:"resourceURI" example:"did:cheqd:testnet:55dbc8bf-fba3-4117-855c-1e0dc1d3bb47/resources/398cee0a-efac-4643-9f4c-74c48c72a14b"`
	CollectionId      string     `json:"resourceCollectionId" example:"55dbc8bf-fba3-4117-855c-1e0dc1d3bb47"`
	ResourceId        string     `json:"resourceId" example:"398cee0a-efac-4643-9f4c-74c48c72a14b"`
	Name              string     `json:"resourceName" example:"Image Resource"`
	ResourceType      string     `json:"resourceType" example:"Image"`
	MediaType         string     `json:"mediaType" example:"image/png"`
	Version           string     `json:"resourceVersion" example:"1"`
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
		Version:           resource.Version,
		Created:           &created,
		Checksum:          resource.Checksum,
		PreviousVersionId: previousVersionId,
		NextVersionId:     nextVersionId,
	}
}

type DereferencedResourceListStruct struct {
	Resources DereferencedResourceList `json:"linkedResourceMetadata,omitempty"`
}

func NewDereferencedResourceListStruct(did string, protoResources []*resourceTypes.Metadata) *DereferencedResourceListStruct {
	resourceList := []DereferencedResource{}
	for _, r := range protoResources {
		resourceList = append(resourceList, *NewDereferencedResource(did, r))
	}

	return &DereferencedResourceListStruct{
		Resources: resourceList,
	}
}

func (e *DereferencedResourceListStruct) AddContext(newProtocol string) {}
func (e *DereferencedResourceListStruct) RemoveContext()                {}
func (e *DereferencedResourceListStruct) GetBytes() []byte              { return []byte{} }

// DereferencedResourceList

type DereferencedResourceList []DereferencedResource

func (e *DereferencedResourceList) AddContext(newProtocol string) {}
func (e *DereferencedResourceList) RemoveContext()                {}
func (e *DereferencedResourceList) GetBytes() []byte              { return []byte{} }

func (e DereferencedResourceList) GetByResourceId(resourceId string) DereferencedResourceList {
	for _, r := range e {
		if r.ResourceId == resourceId {
			return DereferencedResourceList{r}
		}
	}
	return DereferencedResourceList{}
}

func (e DereferencedResourceList) FilterByCollectionId(collectionId string) DereferencedResourceList {
	filteredResources := DereferencedResourceList{}
	for _, r := range e {
		if r.CollectionId == collectionId {
			filteredResources = append(filteredResources, r)
		}
	}
	return filteredResources
}

func (e DereferencedResourceList) FilterByResourceType(resourceType string) DereferencedResourceList {
	filteredResources := DereferencedResourceList{}
	for _, r := range e {
		if r.ResourceType == resourceType {
			filteredResources = append(filteredResources, r)
		}
	}
	return filteredResources
}

func (e DereferencedResourceList) FilterByResourceName(resourceName string) DereferencedResourceList {
	filteredResources := DereferencedResourceList{}
	for _, r := range e {
		if r.Name == resourceName {
			filteredResources = append(filteredResources, r)
		}
	}
	return filteredResources
}

func (e DereferencedResourceList) FilterByVersion(version string) DereferencedResourceList {
	filteredResources := DereferencedResourceList{}
	for _, r := range e {
		if r.Version == version {
			filteredResources = append(filteredResources, r)
		}
	}
	return filteredResources
}

func (e DereferencedResourceList) FindBeforeTime(stime string) (string, error) {
	search_time, err := utils.ParseFromStringTimeToGoTime(stime)
	if err != nil {
		return "", err
	}
	// Firstly - sort versions by Updated time
	versions := e
	sort.Sort(versions)
	if len(versions) == 0 {
		return "", nil
	}
	for _, v := range versions {
		if v.Created.Before(search_time) || v.Created.Equal(search_time) {
			return v.ResourceId, nil
		}
	}
	return "", nil
}

func (e DereferencedResourceList) FindAllBeforeTime(stime string) (DereferencedResourceList, error) {
	l := DereferencedResourceList{}
	search_time, err := utils.ParseFromStringTimeToGoTime(stime)
	if err != nil {
		return l, err
	}
	// Firstly - sort versions by Updated time
	versions := e
	sort.Sort(versions)
	if len(versions) == 0 {
		return l, nil
	}
	for _, v := range versions {
		if v.Created.Before(search_time) || v.Created.Equal(search_time.Add(time.Second)) {
			l = append(l, v)
		}
	}
	return l, nil
}

func (e DereferencedResourceList) AreResourceNamesTheSame() bool {
	if len(e) == 0 {
		return true
	}
	firstName := e[0].Name
	for _, r := range e {
		if r.Name != firstName {
			return false
		}
	}
	return true
}

func (e DereferencedResourceList) AreResourceTypesTheSame() bool {
	if len(e) == 0 {
		return true
	}
	firstType := e[0].ResourceType
	for _, r := range e {
		if r.ResourceType != firstType {
			return false
		}
	}
	return true
}

func (dr DereferencedResourceList) Len() int {
	return len(dr)
}

// Sort in reverse order
func (dr DereferencedResourceList) Less(i, j int) bool {
	return dr[i].Created.After(*dr[j].Created)
}

func (dr DereferencedResourceList) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

// DereferencedResourceData

type DereferencedResourceData []byte

func (e *DereferencedResourceData) AddContext(newProtocol string) {}
func (e *DereferencedResourceData) RemoveContext()                {}
func (e *DereferencedResourceData) GetBytes() []byte              { return []byte(*e) }
