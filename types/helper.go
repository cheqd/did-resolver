package types

import (
	"crypto/sha256"
	"fmt"
)

func FixResourceChecksum(inputChecksum []byte) (hash string) {
	if len(fmt.Sprintf("%x", inputChecksum)) == 64 {
		return fmt.Sprintf("%x", inputChecksum)
	}
	h := sha256.New()
	data := inputChecksum[:len(inputChecksum)-len(h.Sum(nil))]
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}
