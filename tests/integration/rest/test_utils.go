//go:build integration

package rest

import (
	"encoding/json"
	"os"

	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/gomega"
)

type PositiveTestCase struct {
	DidURL                 string
	ResolutionType         string
	EncodingType           string
	ExpectedEncodingType   string
	ExpectedJSONPath       string
	ExpectedStatusCode     int
	ExpectedLocationHeader string
}

type NegativeTestCase struct {
	DidURL             string
	ResolutionType     string
	ExpectedResult     any
	ExpectedStatusCode int
}

type DereferencingResult struct {
	Context               string                         `json:"@context,omitempty"`
	DereferencingMetadata types.DereferencingMetadata    `json:"dereferencingMetadata"`
	ContentStream         *any                           `json:"contentStream"`
	Metadata              types.ResolutionDidDocMetadata `json:"contentMetadata"`
}

type ResourceDereferencingResult struct {
	Context               string                      `json:"@context,omitempty"`
	DereferencingMetadata types.DereferencingMetadata `json:"dereferencingMetadata"`
	ContentStream         *any                        `json:"contentStream"`
	Metadata              *types.DereferencedResource `json:"contentMetadata,omitempty"`
}

func AssertDidDereferencing(expected DereferencingResult, received DereferencingResult) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.DereferencingMetadata.ContentType).To(Equal(received.DereferencingMetadata.ContentType))
	Expect(expected.DereferencingMetadata.ResolutionError).To(Equal(received.DereferencingMetadata.ResolutionError))
	Expect(expected.DereferencingMetadata.DidProperties).To(Equal(received.DereferencingMetadata.DidProperties))
	Expect(expected.ContentStream).To(Equal(received.ContentStream))
	Expect(expected.Metadata).To(Equal(received.Metadata))
}

func AssertDidResolution(expected types.DidResolution, received types.DidResolution) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.ResolutionMetadata.ContentType).To(Equal(received.ResolutionMetadata.ContentType))
	Expect(expected.ResolutionMetadata.ResolutionError).To(Equal(received.ResolutionMetadata.ResolutionError))
	Expect(expected.ResolutionMetadata.DidProperties).To(Equal(received.ResolutionMetadata.DidProperties))
	Expect(expected.Did).To(Equal(received.Did))
	Expect(expected.Metadata).To(Equal(received.Metadata))
}

func AssertResourceDataWithMetadata(expected ResourceDereferencingResult, received ResourceDereferencingResult) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.DereferencingMetadata.ContentType).To(Equal(received.DereferencingMetadata.ContentType))
	Expect(expected.DereferencingMetadata.ResolutionError).To(Equal(received.DereferencingMetadata.ResolutionError))
	Expect(expected.DereferencingMetadata.DidProperties).To(Equal(received.DereferencingMetadata.DidProperties))
	Expect(expected.ContentStream).To(Equal(received.ContentStream))
	Expect(expected.Metadata.ResourceType).To(Equal(received.Metadata.ResourceType))
	Expect(expected.Metadata.ResourceId).To(Equal(received.Metadata.ResourceId))
	Expect(expected.Metadata.Version).To(Equal(received.Metadata.Version))
	Expect(expected.Metadata.Name).To(Equal(received.Metadata.Name))
	Expect(expected.Metadata.CollectionId).To(Equal(received.Metadata.CollectionId))
}

func AssertResourceMetadata(expected ResourceDereferencingResult, received ResourceDereferencingResult) {
	Expect(expected.Context).To(Equal(received.Context))
	Expect(expected.DereferencingMetadata.ContentType).To(Equal(received.DereferencingMetadata.ContentType))
	Expect(expected.DereferencingMetadata.ResolutionError).To(Equal(received.DereferencingMetadata.ResolutionError))
	Expect(expected.DereferencingMetadata.DidProperties).To(Equal(received.DereferencingMetadata.DidProperties))
	Expect(expected.ContentStream).To(Equal(received.ContentStream))
	Expect(expected.Metadata).To(Equal(received.Metadata))
}

func ConvertJsonFileToType(path string, v any) error {
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

func GetTextResource(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(content), nil
}
