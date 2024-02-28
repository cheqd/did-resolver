//go:build integration

package query_test

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

var _ = DescribeTable("Positive: request with common query parameters", func(testCase utils.PositiveTestCase) {
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
		"can get DIDDoc with an existent versionId and versionTime query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"ce298b6f-594b-426e-b431-370d6bc5d3ad",
				"2023-03-06T09:39:49Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/version_id/version_time.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent versionId and transformKeys query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&transformKeys=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"ce298b6f-594b-426e-b431-370d6bc5d3ad",
				types.JsonWebKey2020,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/version_id/transform_key.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc with an existent versionId, versionTime, transformKeys",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=%s&transformKeys=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"ce298b6f-594b-426e-b431-370d6bc5d3ad",
				"2023-03-06T09:39:49Z",
				types.Ed25519VerificationKey2020,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/version_id/version_time_&_transform_key.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)

var _ = DescribeTable("Positive: request with common query parameters (metadata)", func(testCase utils.PositiveTestCase) {
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
		"can get DIDDoc metadata with an existent versionId query parameters and supported value of metadata",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&metadata=true",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"0ce23d04-5b67-4ea6-a315-788588e53f4e",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/metadata/version_id.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent metadata and versionTime query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionTime=%s&metadata=true",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"2023-03-06T09:39:49Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/metadata/version_time.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an existent metadata, versionId, versionTime query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=%s&metadata=true",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				"0ce23d04-5b67-4ea6-a315-788588e53f4e",
				"2023-03-06T09:36:56Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/diddoc_common/metadata/version_id_&_version_time.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)

var (
	serviceId              = "bar"
	expectedLocationHeader = "https://bar.example.com"
)

var _ = DescribeTable("Positive: request with common query parameters (service)", func(testCase utils.PositiveTestCase) {
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
		"can redirect to serviceEndpoint with existent service and versionId query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?service=%s&versionId=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				serviceId,
				"ce298b6f-594b-426e-b431-370d6bc5d3ad",
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: expectedLocationHeader,
		},
	),

	Entry(
		"can redirect to serviceEndpoint with existent service and versionTime query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?service=%s&versionTime=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				serviceId,
				"2023-03-06T09:39:49Z",
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: expectedLocationHeader,
		},
	),

	Entry(
		"can redirect to serviceEndpoint with existent service, relativeRef, versionId, versionTime query parameters",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?service=%s&relativeRef=%s&versionId=%s&versionTime=%s",
				testconstants.TestHostAddress,
				testconstants.SeveralVersionsDID,
				serviceId,
				"\u002Fabout",
				"ce298b6f-594b-426e-b431-370d6bc5d3ad",
				"2023-03-06T09:59:23Z",
			),
			ResolutionType:         testconstants.DefaultResolutionType,
			ExpectedStatusCode:     http.StatusSeeOther,
			ExpectedLocationHeader: fmt.Sprintf("%s/about", expectedLocationHeader),
		},
	),
)
