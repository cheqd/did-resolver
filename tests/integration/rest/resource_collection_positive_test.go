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

var _ = DescribeTable("Positive: get collection of resources", func(testCase positiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing dereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing dereferencingResult
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	assertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get collection of resources with existent DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/collection_of_resources/metadata.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	// TODO: Add test case for getting collection of resources with existent old 16 characters Indy style DID.

	Entry(
		"can get collection of resources with existent old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/collection_of_resources/metadata_32_indy_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
