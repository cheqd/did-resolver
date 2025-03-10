//go:build integration

package diddoc_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get DIDDoc", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		SetHeader("Accept-Encoding", testCase.EncodingType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))
	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))

	var receivedDidResolution types.DidResolution
	unmarshalErr := json.Unmarshal(resp.Body(), &receivedDidResolution)
	if unmarshalErr != nil {
		var resolutionResult types.DidDoc
		var expectedResult types.DidDoc
		Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedResult)).To(BeNil())
		Expect(json.Unmarshal(resp.Body(), &resolutionResult)).To(BeNil())
		utils.AssertDidResolution(types.DidResolution{Did: &expectedResult}, types.DidResolution{Did: &resolutionResult})
	} else {
		var expectedDidResolution types.DidResolution
		Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidResolution)).To(BeNil())
		utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
	}
},

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style mainnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_indy_mainnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent 22 bytes INDY style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleTestnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style mainnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_mainnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent UUID style testnet DID with Chrome Accept header",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:       testconstants.ChromeResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent old 16 characters INDY style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_old_16_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent old 32 characters INDY style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_old_32_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent DID and supported DIDJSON resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:       string(types.DIDJSON),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_did_json.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent DID and supported JSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:       string(types.JSONLD) + ";profile=" + string(types.W3IDDIDRES),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_indy_mainnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent DID and supported DIDJSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:       string(types.DIDJSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_indy_mainnet_did_ld.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent DID and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         "gzip",
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent DID and not supported encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			EncodingType:       testconstants.NotSupportedEncodingType,
			ExpectedJSONPath:   "../../testdata/diddoc/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
