//go:build integration

package query

import (
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get resource", func(testCase utils.PositiveTestCase) {
	var expectedDidDereferencing string
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		SetHeader("Accept-Encoding", testCase.EncodingType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	expectedDidDereferencing, err = utils.GetTextResource(testCase.ExpectedJSONPath)
	Expect(err).To(BeNil())
	Expect(expectedDidDereferencing).To(Equal(string(resp.Body())))

	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
},
	Entry(
		"can get resource with only resourceType (there is only one resource with such name)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/did:cheqd:testnet:b5d70adf-31ca-4662-aa10-d3a54cd8f06c?resourceType=TestType",
			),
			ResolutionType:       testconstants.DefaultResolutionType,
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_data/resource_hello_world.txt",
			ExpectedStatusCode:   http.StatusOK,
		},
	),
)
