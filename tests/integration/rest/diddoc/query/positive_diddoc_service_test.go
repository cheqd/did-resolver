//go:build integration

package query_test

import (
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	ServiceId              = "bar"
	ExpectedLocationHeader = "https://bar.example.com"
)

var _ = DescribeTable("Positive: Get Service param", func(testCase utils.PositiveTestCase) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).NotTo(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))
	Expect(testCase.ExpectedLocationHeader).To(Equal(resp.Header().Get("Location")))
},

	Entry(
		"can redirect to service endpoint",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s",
				SeveralVersionsDID,
				ServiceId,
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: ExpectedLocationHeader,
		},
	),
	Entry(
		"can redirect to service endpoint with relativeRef",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&relativeRef=foo",
				SeveralVersionsDID,
				ServiceId,
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: ExpectedLocationHeader + "foo",
		},
	),
	Entry(
		"can redirect to service endpoint with relativeRef and with versionId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&relativeRef=foo&versionId=%s",
				SeveralVersionsDID,
				ServiceId,
				SeveralVersionVersionId,
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: ExpectedLocationHeader + "foo",
		},
	),

	Entry(
		"can redirect to service endpoint with relativeRef and with versionId and versionTime",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&relativeRef=foo&versionId=%s&versionTime=%s",
				SeveralVersionsDID,
				ServiceId,
				SeveralVersionVersionId,
				SeveralVersionTimeAfter,
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: ExpectedLocationHeader + "foo",
		},
	),
)
