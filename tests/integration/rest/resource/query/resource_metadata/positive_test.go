//go:build integration

package resource_metadata_test

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get Resource Metadata with resourceMetadata query", func() {
	// client := resty.New()

	// resp, err := client.R().
	// 	SetHeader("Accept", testCase.ResolutionType).
	// 	Get(testCase.DidURL)
	// Expect(err).To(BeNil())

	// var receivedResourceData any
	// Expect(json.Unmarshal(resp.Body(), &receivedResourceData)).To(BeNil())
	// Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	// var expectedResourceData any
	// Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedResourceData)).To(BeNil())

	// Expect(expectedResourceData).To(Equal(receivedResourceData))
},

	Entry(
		"can get resource metadata with resourceMetadata=true query parameter",
		// utils.PositiveTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?resourceMetadata=true",
		// 		testconstants.UUIDStyleTestnetDid,
		// 	),
		// 	ResolutionType:     testconstants.DefaultResolutionType,
		// 	ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata.json",
		// 	ExpectedStatusCode: http.StatusOK,
		// },
	),

	Entry(
		"can get resource metadata with resourceMetadata=false query parameter",
		// utils.PositiveTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?resourceMetadata=false",
		// 		testconstants.UUIDStyleTestnetDid,
		// 	),
		// 	ResolutionType:     testconstants.DefaultResolutionType,
		// 	ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata.json",
		// 	ExpectedStatusCode: http.StatusOK,
		// },
	),

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and resourceMetadata query parameter",
		// utils.PositiveTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?resourceMetadata=true",
		// 		testconstants.UUIDStyleTestnetDid,
		// 	),
		// 	ResolutionType:     testconstants.DefaultResolutionType,
		// 	ExpectedJSONPath:   "../../../testdata/query/resource_metadata/metadata_32_indy_did.json",
		// 	ExpectedStatusCode: http.StatusOK,
		// },
	),
)
