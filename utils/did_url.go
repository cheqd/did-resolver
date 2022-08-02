package utils

import (
	"regexp"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
)

var (
	ResourcePath, _ = regexp.Compile(`resources\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	ResourceId, _   = regexp.Compile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
)

func GetResourceId(didUrlPath string) (id string) {
	if !ResourcePath.Match([]byte(didUrlPath)) {
		return ""
	}

	match := ResourceId.FindStringSubmatch(didUrlPath)
	return match[0]
}

func IsDidUrl(didUrl string) bool {
	_, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	return err == nil && (path != "" || query != "" || fragmentId != "")
}
