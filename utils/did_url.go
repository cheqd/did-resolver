package utils

import (
	"regexp"

	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
)

var ResourceId, _ = regexp.Compile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func IsDidUrl(didUrl string) bool {
	_, path, query, fragmentId, err := didUtils.TrySplitDIDUrl(didUrl)
	return err == nil && (path != "" || query != "" || fragmentId != "")
}

func IsValidResourceId(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
