package utils

import (
	"regexp"
)

var ResourcePath, _ = regexp.Compile(`resource\/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func GetResourceId(didUrlPath string) (id string) {
	if !ResourcePath.Match([]byte(didUrlPath)) {
		return ""
	}

	match := ResourcePath.FindStringSubmatch(didUrlPath)
	return match[0]
}