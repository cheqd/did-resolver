//go:build integration

package query_test

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

var _ = DescribeTable("Positive: request with common query parameters", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedResourceData any
	Expect(json.Unmarshal(resp.Body(), &receivedResourceData)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedResourceData any
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedResourceData)).To(BeNil())

	Expect(expectedResourceData).To(Equal(receivedResourceData))
},

	Entry(
		"can get resource with an existent versionId and resourceId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&resourceId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_1.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource with an existent versionTime and resourceId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionTime=%s&resourceId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				"2023-01-25T11:58:11Z",
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_1.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource with an existent versionId, versionTime, resourceCollection, resourceId",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceCollectionId=%s&resourceId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetVersionId,
				"2023-01-25T11:58:11Z",
				testconstants.UUIDStyleTestnetId,
				testconstants.UUIDStyleTestnetDidResourceId,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_1.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource with an existent resourceId and resourceVersionTime",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceId=%s&resourceVersionTime=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetDidResourceId,
				"2023-01-25T12:08:40Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_1.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource with an existent resourceCollectionId, resourceId, resourceName, resourceType, resourceVersion",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceCollectionId=%s&resourceId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s",
				testconstants.TestHostAddress,
				"did:cheqd:testnet:0a5b94d0-a417-48ed-a6f5-4abc9e95888d",
				"0a5b94d0-a417-48ed-a6f5-4abc9e95888d",
				"ef344b53-f2db-44bd-9df3-01259c178704",
				"MuseumPassCredentialSchema",
				"JsonSchemaValidator2018",
				"1.0",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_2.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get resource with an existent resourceCollectionId, resourceId, resourceName, resourceType, resourceVersion, resourceVersionTime",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceCollectionId=%s&resourceId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s",
				testconstants.TestHostAddress,
				"did:cheqd:testnet:0a5b94d0-a417-48ed-a6f5-4abc9e95888d",
				"0a5b94d0-a417-48ed-a6f5-4abc9e95888d",
				"ef344b53-f2db-44bd-9df3-01259c178704",
				"MuseumPassCredentialSchema",
				"JsonSchemaValidator2018",
				"1.0",
				"2023-03-24T11:55:00Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../testdata/query/resource_common/resource_combination_of_queries_2.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
