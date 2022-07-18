package utils

import (
	"regexp"
)

var ResourcePath, _ = regexp.Compile(`resources\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
var ResourceId, _ = regexp.Compile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func GetResourceId(didUrlPath string) (id string) {
	if !ResourcePath.Match([]byte(didUrlPath)) {
		return ""
	}

	match := ResourceId.FindStringSubmatch(didUrlPath)
	return match[0]
}
