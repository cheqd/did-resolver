//go:build integration

package versionTime

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

var SeveralVersionsDID = "did:cheqd:testnet:b5d70adf-31ca-4662-aa10-d3a54cd8f06c"

var _ = DescribeTable("Positive: Get DIDDoc with versionTime query", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidResolution types.DidResolution
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidResolution)).To(BeNil())

	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"can get DIDDoc with versionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				SeveralVersionsDID,
				"2023-03-06T09:39:49.496306968Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an old 16 characters INDY style DID and versionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				"2022-10-13T06:09:04Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_16_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an old 32 characters INDY style DID and versionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				"2022-10-12T08:57:25Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_32_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
