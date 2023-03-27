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

type resourceMetadataPositiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

var _ = DescribeTable("", func(testCase resourceMetadataPositiveTestCase) {
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

		"can get resource metadata with existent DID and resourceId",
		resourceMetadataPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/resource_metadata/metadata.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
