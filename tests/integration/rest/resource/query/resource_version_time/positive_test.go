//go:build integration

package resource_version_time_test

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

var _ = DescribeTable("Positive: Get Collection of Resources with resourceVersionTime query", func(testCase utils.PositiveTestCase) {
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
		"can get collection of resources with an existent resourceVersionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.UUIDStyleTestnetDid,
				"2023-01-25T12:08:40Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/collection_id/metadata_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	// TODO: add unit test for testing get resource with an old 16 characters INDY style DID
	// and resourceVersionTime query parameter.

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceVersionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.OldIndy32CharStyleTestnetDid,
				"2022-10-12T08:57:31Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/collection_id/metadata_32_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
