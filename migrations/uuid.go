package migrations

import (
	"github.com/cheqd/did-resolver/utils"
	"github.com/google/uuid"
)

func MigrateUUIDDid(did string) string {
	method, namespace, id := utils.MustSplitDID(did)
	return utils.JoinDID(method, namespace, MigrateUUIDId(id))
}

func MigrateUUIDId(id string) string {
	// If id is not UUID it should not be changed
	if !utils.IsValidUUID(id) {
		return id
	}

	// If uuid is already normalized, it should not be changed
	if utils.NormalizeUUID(id) == id {
		return id
	}

	newID := uuid.NewSHA1(uuid.Nil, []byte(id))

	return utils.NormalizeUUID(newID.String())
}
