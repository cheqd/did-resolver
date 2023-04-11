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

var _ = DescribeTable("Positive: get resource metadata", func(testCase utils.PositiveTestCase) {
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
		"can get resource metadata with existent DID and resourceId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	// TODO: add test for getting resource metadata with existent old 16 characters Indy style DID
	// and an existent resourceId.

	Entry(
		"can get resource metadata with existent old 32 characters Indy style DID and an existent resourceId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.OldIndy32CharStyleTestnetDidIdentifierResourceId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata_32_indy_did.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource metadata with an existent DID, and supported DIDJSON resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.DIDJSON),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata_did_json.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource metadata with an existent DID, and supported DIDJSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.DIDJSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported JSONLD resolution type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.JSONLD),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         "gzip",
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_metadata/metadata.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc version with an existent DID, and supported gzip encoding type",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			EncodingType:       testconstants.NotSupportedEncodingType,
			ExpectedJSONPath:   "../../testdata/resource_metadata/metadata.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
