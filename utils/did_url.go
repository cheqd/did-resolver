package utils

import (
	"regexp"
)

var ResourcePath, _ = regexp.Compile(`resources\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func GetResourceId(didUrlPath string) (id string) {
	match := ResourcePath.FindStringSubmatch(didUrlPath)
	if len(match) != 1 {
		return ""
	}
	return match[0]
}
