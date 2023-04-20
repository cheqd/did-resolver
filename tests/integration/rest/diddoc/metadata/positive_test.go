//go:build integration

package metadata

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

var _ = DescribeTable("Positive: Get DIDDoc version metadata", func(testCase utils.PositiveTestCase) {
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
		"can get DIDDoc version metadata with an existent 22 bytes INDY style mainnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.IndyStyleMainnetDid,
				testconstants.IndyStyleMainnetVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_indy_mainnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent 22 bytes INDY style testnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.IndyStyleTestnetDid,
				testconstants.IndyStyleTestnetVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent UUID style mainnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleMainnetDid,
				testconstants.UUIDStyleMainnetVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_uuid_mainnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent UUID style testnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent old 16 characters Indy style testnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				"674e6cb5-8d7c-5c50-b0ff-d91bcbcbd5d6",
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_old_16_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent old 32 characters Indy style testnet DID and versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				"1dc202d4-26ee-54a9-b091-8d2e1f609722",
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_old_32_indy_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent DID and versionId, and supported DIDJSON resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:       string(types.DIDJSON),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_did_json.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent DID and versionId, and supported DIDJSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:       string(types.DIDJSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent DID and versionId, and supported JSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:       string(types.JSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent DID and versionId, and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         "gzip",
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version metadata with an existent DID and versionId, and not supported encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			EncodingType:       testconstants.NotSupportedEncodingType,
			ExpectedJSONPath:   "../../testdata/diddoc_version_metadata/diddoc_uuid_testnet_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
