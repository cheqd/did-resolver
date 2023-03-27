//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type resourceCollectionPositiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

var _ = DescribeTable("Positive: get collection of resources", func(testCase resourceCollectionPositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing DereferencingResult
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	Expect(expectedDidDereferencing.Context).To(Equal(receivedDidDereferencing.Context))
	Expect(expectedDidDereferencing.DereferencingMetadata.ContentType).To(Equal(receivedDidDereferencing.DereferencingMetadata.ContentType))
	Expect(expectedDidDereferencing.DereferencingMetadata.ResolutionError).To(Equal(receivedDidDereferencing.DereferencingMetadata.ResolutionError))
	Expect(expectedDidDereferencing.DereferencingMetadata.DidProperties).To(Equal(receivedDidDereferencing.DereferencingMetadata.DidProperties))
	Expect(expectedDidDereferencing.ContentStream).To(Equal(receivedDidDereferencing.ContentStream))
	Expect(expectedDidDereferencing.Metadata).To(Equal(receivedDidDereferencing.Metadata))
},

	Entry(
		"can get collection of resources with existent DID",
		resourceCollectionPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/collection_of_resources/metadata.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
