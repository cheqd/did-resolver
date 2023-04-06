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

var _ = DescribeTable("Positive: Get DIDDoc versions", func(testCase positiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		SetHeader("Accept-Encoding", testCase.encodingType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing dereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing dereferencingResult
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	Expect(testCase.expectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	assertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get DIDDoc versions with an existent DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc versions with an existent old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions_old_16_indy_did.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc versions with an existent old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions_old_32_indy_did.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc versions with an existent DID, and supported DIDJSON resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.UUIDStyleTestnetDid,
			),
			resolutionType:       string(types.DIDJSON),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions_did_json.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported DIDJSONLD resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:       string(types.DIDJSONLD),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported JSONLD resolution type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:       string(types.JSONLD),
			encodingType:         testconstants.DefaultEncodingType,
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported gzip encoding type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:       testconstants.DefaultResolutionType,
			encodingType:         "gzip",
			expectedEncodingType: "gzip",
			expectedJSONPath:     "testdata/diddoc_versions/diddoc_versions.json",
			expectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and not supported encoding type",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			encodingType:       testconstants.NotSupportedEncodingType,
			expectedJSONPath:   "testdata/diddoc_versions/diddoc_versions.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
