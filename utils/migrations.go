package utils

import (
	"crypto/sha256"
	"errors"

	didutils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/google/uuid"
	"github.com/mr-tron/base58"
)

func MigrateIndyStyleDid(did string) string {
	method, namespace, id := didutils.MustSplitDID(did)
	return didutils.JoinDID(method, namespace, MigrateIndyStyleId(id))
}

func MigrateIndyStyleId(id string) string {
	// If id is UUID it should not be changed
	if didutils.IsValidUUID(id) {
		return id
	}

	// Get Hash from current id to make a 32-symbol string
	hash := sha256.Sum256([]byte(id))
	// Indy-style identifier is 16-byte base58 string
	return base58.Encode(hash[:16])
}

func MigrateUUIDDid(did string) string {
	method, namespace, id := didutils.MustSplitDID(did)
	return didutils.JoinDID(method, namespace, MigrateUUIDId(id))
}

func MigrateUUIDId(id string) string {
	// If id is not UUID it should not be changed
	if !didutils.IsValidUUID(id) {
		return id
	}

	// If uuid is already normalized, it should not be changed
	if didutils.NormalizeUUID(id) == id {
		return id
	}

	newId := uuid.NewSHA1(uuid.Nil, []byte(id))
	return didutils.NormalizeUUID(newId.String())
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
