//go:build integration

package rest

import (
	"encoding/json"
	"os"

	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/gomega"
)

type positiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

type negativeTestCase struct {
	didURL             string
	resolutionType     string
	expectedResult     any
	expectedStatusCode int
}

type dereferencingResult struct {
	Context               string                         `json:"@context,omitempty"`
	DereferencingMetadata types.DereferencingMetadata    `json:"dereferencingMetadata"`
	ContentStream         *any                           `json:"contentStream"`
	Metadata              types.ResolutionDidDocMetadata `json:"contentMetadata"`
}

func assertDidDereferencing(expected dereferencingResult, received dereferencingResult) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.DereferencingMetadata.ContentType).To(Equal(received.DereferencingMetadata.ContentType))
	Expect(expected.DereferencingMetadata.ResolutionError).To(Equal(received.DereferencingMetadata.ResolutionError))
	Expect(expected.DereferencingMetadata.DidProperties).To(Equal(received.DereferencingMetadata.DidProperties))
	Expect(expected.ContentStream).To(Equal(received.ContentStream))
	Expect(expected.Metadata).To(Equal(received.Metadata))
}

func assertDidResolution(expected types.DidResolution, received types.DidResolution) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.ResolutionMetadata.ContentType).To(Equal(received.ResolutionMetadata.ContentType))
	Expect(expected.ResolutionMetadata.ResolutionError).To(Equal(received.ResolutionMetadata.ResolutionError))
	Expect(expected.ResolutionMetadata.DidProperties).To(Equal(received.ResolutionMetadata.DidProperties))
	Expect(expected.Did).To(Equal(received.Did))
	Expect(expected.Metadata).To(Equal(received.Metadata))
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
