//go:build integration

package data_test

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

var _ = DescribeTable("Positive: Get resource data with metadata", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		SetHeader("Accept-Encoding", testCase.EncodingType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedResourceDataWithMetadata utils.ResourceDereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedResourceDataWithMetadata)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedResourceDataWithMetadata utils.ResourceDereferencingResult
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedResourceDataWithMetadata)).To(BeNil())

	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	utils.AssertResourceDataWithMetadata(expectedResourceDataWithMetadata, receivedResourceDataWithMetadata)
},

	Entry(
		"can get resource with metadata with an existent DID, and supported JSONLD resolution type and W3DIDUrl dereferencing param",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.JSONLD) + ";profile=" + string(types.W3IDDIDURL),
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_data_with_metadata/resource.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource with metadata with an existent DID, with multiple headers without q values",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.JSONLD) + ";profile=" + string(types.W3IDDIDURL) + ",*/*",
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_data_with_metadata/resource.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),

	Entry(
		"can get resource with metadata with an existent DID, with multiple header and q values",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:       string(types.JSONLD) + ";profile=" + string(types.W3IDDIDURL) + ";q=1.0,application/json;q=0.9,image/png;q=0.7",
			EncodingType:         testconstants.DefaultEncodingType,
			ExpectedEncodingType: "gzip",
			ExpectedJSONPath:     "../../testdata/resource_data_with_metadata/resource.json",
			ExpectedStatusCode:   http.StatusOK,
		},
	),
)
