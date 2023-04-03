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

var _ = DescribeTable("Positive: Get DIDDoc", func(testCase positiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidResolution types.DidResolution
	Expect(convertJsonFileToType(testCase.expectedJSONPath, &expectedDidResolution)).To(BeNil())

	assertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style mainnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.IndyStyleMainnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_indy_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style testnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.IndyStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style mainnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.UUIDStyleMainnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_uuid_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style testnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.UUIDStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_uuid_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent old 16 characters INDY style testnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_old_16_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent old 32 characters INDY style testnet DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_old_32_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)