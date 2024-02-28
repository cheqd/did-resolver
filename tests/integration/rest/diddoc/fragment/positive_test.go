//go:build integration

package fragment_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get DID#fragment", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		SetHeader("Accept-Encoding", testCase.EncodingType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing utils.DereferencingResult
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get verificationMethod section with an existent DID#fragment",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent old 16 characters Indy style DID#fragment",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey-1",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_old_16_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent old 32 characters Indy style DID#fragment",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey-1",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_old_32_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported DIDJSON resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.DIDJSON),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment_did_json.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported DIDJSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.DIDJSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported JSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.JSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment_json_ld.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported default encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment_json_ld.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         "gzip",
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment_json_ld.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get verificationMethod section with an existent DID#fragment and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%skey1",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.NotSupportedEncodingType,
			ExpectedEncodingType: "",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/verification_method_did_fragment_json_ld.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/service_endpoint_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment and supported DIDJSON resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.DIDJSON),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/service_endpoint_did_fragment_did_json.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment and supported DIDJSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.DIDJSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/service_endpoint_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment and supported JSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       string(types.JSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/service_endpoint_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         "gzip",
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_fragment/service_endpoint_did_fragment.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get serviceEndpoint section with an existent DID#fragment and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%swebsite",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			EncodingType:       testconstants.NotSupportedEncodingType,
			ExpectedJSONPath:   "../../testdata/diddoc_fragment/service_endpoint_did_fragment.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
