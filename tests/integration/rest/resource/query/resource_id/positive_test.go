//go:build integration

package resource_id_test

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get Resource with resourceId query", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedResourceData any
	Expect(json.Unmarshal(resp.Body(), &receivedResourceData)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedResourceData any
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedResourceData)).To(BeNil())

	Expect(expectedResourceData).To(Equal(receivedResourceData))
},

	Entry(
		"can get resource with an existent resourceId query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_id/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	// TODO: add unit test for testing get resource with an old 16 characters INDY style DID
	// and resourceId query parameter.

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and an existent resourceId query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceId=%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.OldIndy32CharStyleTestnetDidIdentifierResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_id/resource_32_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
