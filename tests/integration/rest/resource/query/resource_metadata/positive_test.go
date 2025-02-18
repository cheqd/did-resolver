//go:build integration

package resource_metadata_test

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

var _ = DescribeTable("Positive: Get Resource Metadata with resourceMetadata query", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing utils.DereferencingResult
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get resource metadata with resourceMetadata=true query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	// Entry(
	// 	"can get resource metadata with resourceMetadata=false query parameter",
	// 	utils.PositiveTestCase{
	// 		DidURL: fmt.Sprintf(
	// 			"http://%s/1.0/identifiers/%s?resourceMetadata=false",
	//			testconstants.TestHostAddress,
	// 			testconstants.UUIDStyleTestnetDid,
	// 		),
	// 		ResolutionType:     testconstants.DefaultResolutionType,
	// 		ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata.json",
	// 		ExpectedStatusCode: http.StatusOK,
	// 	},
	// ),

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and resourceMetadata query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_32_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get filtered list of resources in didDocumentMetadata when acceptHeader is W3IDDIDRES",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceName=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.ExistentResourceName,
			),
			ResolutionType:     string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_did_res.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
