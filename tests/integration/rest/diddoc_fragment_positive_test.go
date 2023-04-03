//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	testconstants "github.com/cheqd/did-resolver/tests/constants"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get DID#fragment", func(testCase positiveTestCase) {
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
		"can get verificationMethod section with an existent DID#fragment",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey1",
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/verification_method_fragment.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent old 16 characters Indy style DID#fragment",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.OldIndy16CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/verification_method_old_16_did_fragment.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent old 32 characters Indy style DID#fragment",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.OldIndy32CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/verification_method_old_32_did_fragment.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%swebsite",
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/service_endpoint_fragment.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
