//go:build integration

package resource_version_test

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

var severalResourcesDID = "did:cheqd:testnet:0a5b94d0-a417-48ed-a6f5-4abc9e95888d"

var _ = DescribeTable("Positive: Get Resource with resourceVersion query", func(testCase utils.PositiveTestCase) {
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
		"can get resource with an existent resourceVersion query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersion=1.0",
				testconstants.SUTHost,
				severalResourcesDID,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	// TODO: add unit test for testing get resource with:
	// - an old 16 characters INDY style DID and resourceVersion query parameter
	// - an old 32 characters INDY style DID and resourceVersion query parameter
)
