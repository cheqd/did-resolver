package utils

import (
	"strings"

	migrations "github.com/cheqd/cheqd-node/app/migrations/helpers"
	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
)

func IsValidResourceId(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func IsValidIndyV1ID(id string) bool {
	return len(id) == 16 && didUtils.IsValidBase58(id) ||
		len(id) == 32 && didUtils.IsValidBase58(id)
}

func IsValidUUIDV1ID(id string) bool {
	return didUtils.IsValidUUID(id) && strings.ToLower(id) != id
}

func IsMigrationNeeded(id string) bool {
	return IsValidIndyV1ID(id) || IsValidUUIDV1ID(id)

}

func MigrateDID(did string) string {
	did = migrations.MigrateIndyStyleDid(did)
	did = migrations.MigrateUUIDDid(did)

	return did
}
