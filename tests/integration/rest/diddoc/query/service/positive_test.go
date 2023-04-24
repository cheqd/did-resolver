//go:build integration

package service

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
	SeveralVersionsDID     = "did:cheqd:testnet:b5d70adf-31ca-4662-aa10-d3a54cd8f06c"
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
		"can redirect to serviceEndpoint with an existent service query parameter",
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
		"can redirect to serviceEndpoint with an existent service and a valid relativeRef URI query parameters",
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

	// TODO: add unit test for testing:
	// - old 16 characters INDY style DID with service query parameter
	// - old 32 characters INDY style DID with service query parameter
)
