package types

import (
	"sort"
	"time"

	"github.com/cheqd/did-resolver/utils"
)

type DidDocMetadataList []ResolutionDidDocMetadata

func (dd DidDocMetadataList) GetByVersionId(versionId string) DidDocMetadataList {
	for _, r := range dd {
		if r.VersionId == versionId {
			return DidDocMetadataList{r}
		}
	}
	return DidDocMetadataList{}
}

// Returns VersionId if there is a version before the given time
// Otherwise NotFound error
func (e DidDocMetadataList) FindActiveForTime(stime string) (string, error) {
	search_time, err := utils.ParseFromStringTimeToGoTime(stime)
	if err != nil {
		return "", err
	}
	// Firstly - sort versions by Updated time
	versions := e
	sort.Sort(versions)
	for _, version := range versions {
		if version.Updated != nil && (version.Updated.Before(search_time) || version.Updated.Equal(search_time)) {
			return version.VersionId, nil
		}
		if version.Updated == nil && (version.Created.Before(search_time) || version.Created.Equal(search_time)) {
			return version.VersionId, nil
		}
	}
	return "", nil
}

func (dd DidDocMetadataList) GetResourcesBeforeNextVersion(versionId string) DereferencedResourceList {
	if len(dd) == 1 {
		return dd[0].Resources
	}
	// If versionId == the latest versionId then return all the resources
	if len(dd) > 1 && dd[0].VersionId == versionId {
		return dd[0].Resources
	}

	var previous ResolutionDidDocMetadata = dd[0]
	for _, r := range dd[1:] {
		if r.VersionId == versionId {
			var timeBefore string
			if previous.Updated != nil {
				timeBefore = previous.Updated.Format(time.RFC3339Nano)
			} else {
				timeBefore = previous.Created.Format(time.RFC3339Nano)
			}
			filteredResources, err := previous.Resources.FindAllBeforeTime(timeBefore)
			if err != nil {
				return DereferencedResourceList{}
			}
			return filteredResources
		}
		previous = r
	}
	return DereferencedResourceList{}
}

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

// Interface implementation

func (d DidDocMetadataList) GetContentType() string { return "" }

func (d DidDocMetadataList) GetBytes() []byte { return []byte{} }

func (d DidDocMetadataList) IsRedirect() bool { return false }
