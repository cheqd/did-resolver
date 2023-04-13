// go:build integration

package query

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	testconstants "github.com/cheqd/did-resolver/tests/constants"
// 	utils "github.com/cheqd/did-resolver/tests/integration/rest"

// 	"github.com/go-resty/resty/v2"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// var _ = DescribeTable("Positive: Get resource", func(testCase utils.PositiveTestCase) {
// 	client := resty.New()
// 	fmt.Println(testCase.DidURL)

// 	resp, err := client.R().
// 		SetHeader("Accept", testCase.ResolutionType).
// 		SetHeader("Accept-Encoding", testCase.EncodingType).
// 		Get(testCase.DidURL)
// 	Expect(err).To(BeNil())

// 	var receivedDidDereferencing utils.DereferencingResult
// 	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
// 	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

// 	var expectedDidDereferencing utils.DereferencingResult
// 	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidDereferencing)).To(BeNil())

// 	Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
// 	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
// },

// 	Entry(
// 		"can get resource with only resourceId",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceId=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),

// 	Entry(
// 		"can get resource with only resourceName (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with resourceVersionTime",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceVersionTime=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				testconstants.ExistentResourceVersionTimeAfter,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with combination resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with combination resourceId, resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with combination resourceId, resourceVersionTime, resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s&resourceVersionTime=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 				testconstants.ExistentResourceVersionTimeAfter,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with combination resourceId, resourceVersionTime, resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s&resourceVersionTime=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 				testconstants.ExistentResourceVersionTimeAfter,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	Entry(
// 		"can get resource with combination versionId, resourceId, resourceVersionTime, resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s&resourceVersionTime=%s&versionId=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 				testconstants.ExistentResourceVersionTimeAfter,
// 				testconstants.UUIDStyleTestnetVersionId,
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),

// 	Entry(
// 		"can get resource with combination versionTime, resourceId, resourceVersionTime, resourceName and resourceType (there is only one resource with such name)",
// 		utils.PositiveTestCase{
// 			DidURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s&resourceVersionTime=%s&versionTime=%s",
// 				testconstants.UUIDStyleTestnetDid,
// 				strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 				testconstants.ExistentResourceType,
// 				testconstants.UUIDStyleTestnetDidResourceId,
// 				testconstants.ExistentResourceVersionTimeAfter,
// 				"2023-01-26T11:58:10.39Z",
// 			),
// 			ResolutionType:       testconstants.DefaultResolutionType,
// 			EncodingType:         testconstants.DefaultEncodingType,
// 			ExpectedEncodingType: "gzip",
// 			ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 			ExpectedStatusCode:   http.StatusOK,
// 		},
// 	),
// 	// ToDo uncomment this test after allowing versionId and versionTime in the same request

// 	// Entry(
// 	// 	"can get resource with combination versionId, versionTime, resourceId, resourceVersionTime, resourceName and resourceType (there is only one resource with such name)",
// 	// 	utils.PositiveTestCase{
// 	// 		DidURL: fmt.Sprintf(
// 	// 			"http://localhost:8080/1.0/identifiers/%s?resourceName=%s&resourceType=%s&resourceId=%s&resourceVersionTime=%s&versionTime=%s&versionId=%s",
// 	// 			testconstants.UUIDStyleTestnetDid,
// 	// 			strings.ReplaceAll(testconstants.ExistentResourceName, " ", "%20"),
// 	// 			testconstants.ExistentResourceType,
// 	// 			testconstants.UUIDStyleTestnetDidResourceId,
// 	// 			testconstants.ExistentResourceVersionTimeAfter,
// 	// 			"2023-01-26T11:58:10.39Z",
// 	// 			testconstants.UUIDStyleTestnetVersionId,
// 	// 		),
// 	// 		ResolutionType:       testconstants.DefaultResolutionType,
// 	// 		EncodingType:         testconstants.DefaultEncodingType,
// 	// 		ExpectedEncodingType: "gzip",
// 	// 		ExpectedJSONPath:     "../../testdata/resource_data/resource.json",
// 	// 		ExpectedStatusCode:   http.StatusOK,
// 	// 	},
// 	// ),
// )



	