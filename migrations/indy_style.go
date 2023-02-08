package migrations

import (
	"crypto/sha256"

	"github.com/cheqd/did-resolver/utils"
	"github.com/mr-tron/base58"
)

func MigrateIndyStyleDid(did string) string {
	method, namespace, id := utils.MustSplitDID(did)
	return utils.JoinDID(method, namespace, MigrateIndyStyleID(id))
}

func MigrateIndyStyleID(id string) string {
	// If id is UUID it should not be changed
	if utils.IsValidUUID(id) {
		return id
	}

	// Get Hash from current id to make a 32-symbol string
	hash := sha256.Sum256([]byte(id))
	// Indy-style identifier is 16-byte base58 string
	return base58.Encode(hash[:16])
}
