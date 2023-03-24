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

type getDidDocPositiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

var _ = DescribeTable("Positive: Get DIDDoc", func(testCase getDidDocPositiveTestCase) {
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

	Expect(expectedDidResolution.Context).To(Equal(receivedDidResolution.Context))
	Expect(expectedDidResolution.ResolutionMetadata.ContentType).To(Equal(receivedDidResolution.ResolutionMetadata.ContentType))
	Expect(expectedDidResolution.ResolutionMetadata.ResolutionError).To(Equal(receivedDidResolution.ResolutionMetadata.ResolutionError))
	Expect(expectedDidResolution.ResolutionMetadata.DidProperties).To(Equal(receivedDidResolution.ResolutionMetadata.DidProperties))
	Expect(expectedDidResolution.Did).To(Equal(receivedDidResolution.Did))
	Expect(expectedDidResolution.Metadata).To(Equal(receivedDidResolution.Metadata))
},

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style mainnet DID",
		getDidDocPositiveTestCase{
			didURL:             fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.IndyStyleMainnetDid),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_indy_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style testnet DID",
		getDidDocPositiveTestCase{
			didURL:             fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.IndyStyleTestnetDid),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style mainnet DID",
		getDidDocPositiveTestCase{
			didURL:             fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.UUIDStyleMainnetDid),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_uuid_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style testnet DID",
		getDidDocPositiveTestCase{
			didURL:             fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.UUIDStyleTestnetDid),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc/diddoc_uuid_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
