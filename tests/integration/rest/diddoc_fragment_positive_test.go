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

type getDidDocFragmentPositiveTestCase struct {
	didURL             string
	resolutionType     string
	expectedJSONPath   string
	expectedStatusCode int
}

var _ = DescribeTable("Positive: Get DID#fragment", func(testCase getDidDocFragmentPositiveTestCase) {
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
		"can get DIDDoc with an existent 22 bytes INDY style mainnet DID#fragment",
		getDidDocFragmentPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey1",
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/diddoc_fragment_indy_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style testnet DID#fragment",
		getDidDocFragmentPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.IndyStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/diddoc_fragment_indy_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style mainnet DID#fragment",
		getDidDocFragmentPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.UUIDStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/diddoc_fragment_uuid_mainnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style testnet DID#fragment",
		getDidDocFragmentPositiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.UUIDStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedJSONPath:   "testdata/diddoc_fragment/diddoc_fragment_uuid_testnet_did.json",
			expectedStatusCode: http.StatusOK,
		},
	),
)
