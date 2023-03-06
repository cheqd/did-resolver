package utils

import (
	"errors"

	"github.com/mr-tron/base58"
)

const (
	IndyIDLength = 16
)

func IsValidIndyID(data string) bool {
	bytes, err := base58.Decode(data)
	if err != nil {
		return false
	}
	return len(bytes) == IndyIDLength
}

func ValidateID(id string) error {
	isValidID := IsValidIndyID(id) || IsValidUUID(id)

	if !isValidID {
		return errors.New("unique id should be one of: 16 bytes of decoded base58 string or UUID")
	}

	return nil
}

func IsValidID(id string) bool {
	err := ValidateID(id)
	return err == nil
}
