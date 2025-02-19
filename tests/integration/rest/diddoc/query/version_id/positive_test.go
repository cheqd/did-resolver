//go:build integration

package versionId

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

var _ = DescribeTable("Positive: Get DIDDoc with versionId query", func(testCase utils.PositiveTestCase) {
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
		"can get DIDDoc with versionId query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"0ce23d04-5b67-4ea6-a315-788588e53f4e",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/diddoc/diddoc_version_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionId query parameter with Chrome Accept header",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"0ce23d04-5b67-4ea6-a315-788588e53f4e",
			),
			ResolutionType:     testconstants.ChromeResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/diddoc/diddoc_version_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an old 16 characters INDY style DID and versionId query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid,
				"674e6cb5-8d7c-5c50-b0ff-d91bcbcbd5d6",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/diddoc/diddoc_version_16_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an old 32 characters INDY style DID and versionId query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				"1dc202d4-26ee-54a9-b091-8d2e1f609722",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/diddoc/diddoc_version_32_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
