//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get resource data", func(testCase positiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		SetHeader("Accept-Encoding", testCase.encodingType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedResourceData any
	Expect(json.Unmarshal(resp.Body(), &receivedResourceData)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedResourceData any
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedResourceData)).To(BeNil())

	Expect(testCase.expectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	Expect(expectedResourceData).To(Equal(receivedResourceData))
},

	Entry(
		"can get resource data with an existent DID and existent resourceId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	// TODO: Add test for getting resource data with an existent old 16 characters Indy style DID
	// and existent resourceId.

	Entry(
		"can get resource data with an existent old 32 characters Indy style DID and existent resourceId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource_32_indy_did.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource data with an existent DID, and supported DIDJSON resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:       string(types.DIDJSON),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource data with an existent DID, and supported DIDJSONLD resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:       string(types.DIDJSONLD),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource data with an existent DID, and supported JSONLD resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:       string(types.JSONLD),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource data with an existent DID, and supported gzip encoding type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         "gzip",
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/resource_data/resource.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource data with an existent DID, and not supported encoding type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleTestnetDid,
				"9ba3922e-d5f5-4f53-b265-fc0d4e988c77",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			encodingType:       testconstants.NotSupportedEncodingType,
			expectedJSONPath:   "testdata/resource_data/resource.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
