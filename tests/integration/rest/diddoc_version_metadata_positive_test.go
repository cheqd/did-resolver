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

var _ = DescribeTable("Positive: Get DIDDoc metadata for a specific version of DIDDoc", func(testCase positiveTestCase) {
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
		"can get DIDDoc metadata with an existent 22 bytes INDY style mainnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.IndyStyleMainnetDid,
				"4fa8e367-c70e-533e-babf-3732d9761061",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_indy_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent 22 bytes INDY style testnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.IndyStyleTestnetDid,
				"60bb3b62-e0f0-545b-a552-63aab5cd1aef",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent UUID style mainnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleMainnetDid,
				"76e546ee-78cd-5372-b34e-8b47461626e1",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_uuid_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent UUID style testnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				"e5615fc2-6f13-42b1-989c-49576a574cef",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent old 16 characters Indy style testnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				"674e6cb5-8d7c-5c50-b0ff-d91bcbcbd5d6",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_old_16_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent old 32 characters Indy style testnet DID and versionId",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				"1dc202d4-26ee-54a9-b091-8d2e1f609722",
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_version_metadata/diddoc_old_32_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
