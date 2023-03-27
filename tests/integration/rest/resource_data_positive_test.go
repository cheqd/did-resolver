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

type getResourceDataPositiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

var _ = DescribeTable("Positive: Get resource data", func(testCase getResourceDataPositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedResourceData any
	Expect(json.Unmarshal(resp.Body(), &receivedResourceData)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedResourceData any
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedResourceData)).To(BeNil())
	Expect(expectedResourceData).To(Equal(receivedResourceData))
},

	Entry(
		"can get resource data with an existent mainnet DID and existent resourceId",
		getResourceDataPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/resource_data/resource.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
