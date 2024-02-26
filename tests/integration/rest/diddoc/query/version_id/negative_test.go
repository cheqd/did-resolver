//go:build integration

package versionId

import (
	"encoding/json"
	"fmt"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
	errors "github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get DIDDoc with versionId query", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var receivedDidDereferencing types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())

	expectedDidDereferencing := testCase.ExpectedResult.(types.DidResolution)
	utils.AssertDidResolution(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"cannot get DIDDoc with not existent versionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid versionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent versionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=2023-03-06T09:59:21Z",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.SeveralVersionsDIDVersionId,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not supported transformKeys value query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&transformKeys=EDDSA",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.SeveralVersionsDIDVersionId,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&service=%s",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.SeveralVersionsDIDVersionId,
				testconstants.NotExistentService,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent versionTime and service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?versionId=%s&versionTime=%s&service=%s",
				testconstants.SUTHost,
				testconstants.SeveralVersionsDID,
				testconstants.SeveralVersionsDIDVersionId,
				"2023-03-06T09:59:21Z",
				testconstants.NotExistentService,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: testconstants.SeveralVersionsDIDIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),
)
