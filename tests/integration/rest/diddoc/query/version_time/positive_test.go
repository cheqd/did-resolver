//go:build integration

package versionTime

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

	Entry(
		"can get DIDDoc with versionTime query parameter (Layout format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("03/06 09:39:50AM '23 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (ANSIC format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Mon Mar 06 09:39:50 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (UnixDate format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Mon Mar 06 09:39:50 UTC 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RubyDate format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Mon Mar 06 09:39:50 +0000 2023"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC822 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("06 Mar 23 09:40 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC822Z format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("06 Mar 23 09:40 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC850 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Monday, 06-Mar-23 09:39:50 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC1123 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Mon, 06 Mar 2023 09:39:50 UTC"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC1123Z format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("Mon, 06 Mar 2023 09:39:50 +0000"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC3339 format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				"2023-03-06T09:39:50Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (RFC3339Nano format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				"2023-03-06T09:39:49.496306968Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (DateTime format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				url.QueryEscape("2023-03-06 09:39:50"),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with versionTime query parameter (DateOnly format)",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				testconstants.SeveralVersionsDID,
				"2023-03-07",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/version_time/diddoc_version_time_date_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
