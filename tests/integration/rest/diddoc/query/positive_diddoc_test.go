//go:build integration

package query

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

var (
	SeveralVersionsDID        = "did:cheqd:testnet:b5d70adf-31ca-4662-aa10-d3a54cd8f06c"
	SeveralVersionTimeAfter   = "2023-03-06T10:59:22.04Z"
	SeveralVersionTimeBetween = "2023-03-06T09:45:22.04Z"
	SeveralVersionTimeBefore  = "2023-03-06T08:59:22.04Z"
	SeveralVersionVersionId   = "f790c9b9-4817-4b31-be43-b198e6e18071"
)

var _ = DescribeTable("Positive: Get DIDDoc", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		SetHeader("Accept-Encoding", testCase.EncodingType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidResolution types.DidResolution
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidResolution)).To(BeNil())

	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"VersionId: can get DIDDoc with an existent UUID style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionId=%s",
				SeveralVersionsDID,
				SeveralVersionVersionId,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_several_versions.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),
	Entry(
		"VersionTime: can get DIDDoc with an existent UUID style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				SeveralVersionsDID,
				SeveralVersionTimeAfter,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_several_versions.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),
	Entry(
		"VersionTime Between: can get DIDDoc with an existent UUID style testnet DID",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				SeveralVersionsDID,
				SeveralVersionTimeBetween,
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/diddoc/diddoc_uuid_testnet_several_versions_between.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),
)
