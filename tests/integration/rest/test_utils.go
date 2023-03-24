package rest

import (
	"encoding/json"
	"os"

	"github.com/cheqd/did-resolver/types"
)

type DereferencingResult struct {
	Context               string                         `json:"@context,omitempty"`
	DereferencingMetadata types.DereferencingMetadata    `json:"dereferencingMetadata"`
	ContentStream         *any                           `json:"contentStream"`
	Metadata              types.ResolutionDidDocMetadata `json:"contentMetadata"`
}

func convertJsonFileToType(path string, v any) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&v); err != nil {
		return err
	}

	return nil
}
