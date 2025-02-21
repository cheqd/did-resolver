//go:build integration

package resource_metadata_test

import (
	"encoding/json"
	"fmt"
	"github.com/cheqd/did-resolver/types"
	"net/http"

	"github.com/cheqd/did-resolver/types"

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

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidResolution types.DidResolution
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidResolution)).To(BeNil())

	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"can get resource metadata with resourceMetadata=true query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:     string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource metadata with resourceMetadata=false query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=false",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:     string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_false.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and resourceMetadata query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_32_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
