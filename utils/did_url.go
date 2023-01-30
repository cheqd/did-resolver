package utils

import (
	"errors"
	"regexp"

	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
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

func ValidateV1ID(id string) error {
	isValidId := len(id) == 16 && didutils.IsValidBase58(id) ||
		len(id) == 32 && didutils.IsValidBase58(id) ||
		didutils.IsValidUUID(id)

	if !isValidId {
		return errors.New("unique id should be one of: 16 symbols base58 string, 32 symbols base58 string, or UUID")
	}

	return nil
}

func IsValidV1ID(id string) bool {
	err := ValidateV1ID(id)
	return err == nil
}
