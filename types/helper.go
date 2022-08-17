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

func AddElemToSet(set []string, newElement string) []string {
	if set == nil {
		set = []string{}
	}
	for _, c := range set {
		if c == newElement {
			return set
		}
	}
	return append(set, newElement)
}
