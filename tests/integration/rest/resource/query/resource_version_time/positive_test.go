//go:build integration

package resource_version_time_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"

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

	// TODO: add unit test for testing get resource with an old 16 characters INDY style DID
	// and resourceVersionTime query parameter.

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceVersionTime query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				"2022-10-12T08:57:31Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource_32_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (Layout format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("01/25 00:08:40PM '23 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (ANSIC format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wed Jan 25 12:08:40 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (UnixDate format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wed Jan 25 12:08:40 UTC 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RubyDate format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wed Jan 25 12:08:40 +0000 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC822 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("25 Jan 23 12:09 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC822Z format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("25 Jan 23 12:09 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC850 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wednesday, 25-Jan-23 12:08:40 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC1123 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wed, 25 Jan 2023 12:08:40 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC1123Z format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("Wed, 25 Jan 2023 12:08:40 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC3339 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				"2023-01-25T12:08:40Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (RFC3339Nano format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				"2023-01-25T12:08:40.0Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (DateTime format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				url.QueryEscape("2023-01-25 12:08:40"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter (DateOnly format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				"2023-01-26",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
	Entry(
		"can get resource metadata with resourceMetadata=true query parameter and dereferencing profile",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceMetadata=true",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType:     string(types.JSONLD) + ";profile=" + types.W3IDDIDURL,
			ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_did_res.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
	Entry(
		"returns resources sorted in descending order by creation time with resourceVersionTime query",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				"did:cheqd:testnet:4bfaac0c-4cfc-44af-8aa7-577e05a630b5",
				"2025-06-02T23:59:59Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/resource_version_time/resource_sorted_desc.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
