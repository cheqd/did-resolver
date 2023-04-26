//go:build integration

package versionTime

import (
	"encoding/json"
	"fmt"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var SeveralVersionsDIDIdentifier = "b5d70adf-31ca-4662-aa10-d3a54cd8f06c"

var _ = DescribeTable("Negative: Get DIDDoc with versionTime query", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())

	expectedDidResolution := testCase.ExpectedResult.(types.DidResolution)
	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"cannot get DIDDoc with not existent versionTime query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s",
				SeveralVersionsDID,
				"2023-03-06T09:36:56.56204903Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        SeveralVersionsDID,
						MethodSpecificId: SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid versionTime query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=2023-03-06",
				SeveralVersionsDID,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        SeveralVersionsDID,
						MethodSpecificId: SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionTime, but not existent versionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s&versionId=%s",
				SeveralVersionsDID,
				"2023-03-06T09:39:49.496306968Z",
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        SeveralVersionsDID,
						MethodSpecificId: SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionTime, but not existent service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s&service=notexistent",
				SeveralVersionsDID,
				"2023-03-06T09:39:49.496306968Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        SeveralVersionsDID,
						MethodSpecificId: SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.NotFoundHttpCode,
		},
	),

	Entry(

		"cannot get DIDDoc with an existent versionTime, but not existent versionId and service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?versionTime=%s&versionId=%s&service=notexistent",
				SeveralVersionsDID,
				"2023-03-06T09:39:49.496306968Z",
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        SeveralVersionsDID,
						MethodSpecificId: SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.NotFoundHttpCode,
		},
	),
)
