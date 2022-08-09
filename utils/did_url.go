package utils

import (
	"regexp"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var (
	ResourceId, _                      = regexp.Compile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	ResourceDataPath, _                = regexp.Compile(`\/resources\/` + ResourceId.String() + `$`)
	ResourceHeaderPath, _              = regexp.Compile(`\/resources\/` + ResourceId.String() + `\/metadata$`)
	CollectionResourcesPath, _         = regexp.Compile(`\/resources\/all$`)
	CollectionResourcesPathRedirect, _ = regexp.Compile(`\/resources\/$`)
)

func GetResourceId(didUrlPath string) (id string) {
	match := ResourceId.FindStringSubmatch(didUrlPath)
	if len(match) != 1 {
		return ""
	}
	return match[0]
}

func IsResourceDataPath(didUrlPath string) bool {
	return ResourceDataPath.MatchString(didUrlPath)
}

func IsResourceHeaderPath(didUrlPath string) bool {
	return ResourceHeaderPath.MatchString(didUrlPath)
}

func IsCollectionResourcesPath(didUrlPath string) bool {
	return CollectionResourcesPath.MatchString(didUrlPath)
}

func IsCollectionResourcesPathRedirect(didUrlPath string) bool {
	return CollectionResourcesPathRedirect.MatchString(didUrlPath)
}

func IsDidUrl(didUrl string) bool {
	_, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	return err == nil && (path != "" || query != "" || fragmentId != "")
}
