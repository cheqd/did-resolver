//go:build integration

package query_test

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

var _ = DescribeTable("Negative: request with common query parameters", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	expectedDidDereferencing := testCase.ExpectedResult.(utils.DereferencingResult)
	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"cannot get DIDDoc with combination of metadata and relativeRef query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&relativeRef=\u002Fabout",
				testconstants.SeveralVersionsDID,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceId=%s",
				testconstants.SeveralVersionsDID,
				"5e16a3f9-7c6e-4b6b-8e28-20f56780ee25",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceName query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceName=%s",
				testconstants.SeveralVersionsDID,
				"TestResource",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceType query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceType=%s",
				testconstants.SeveralVersionsDID,
				"TestType",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceVersionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceVersionTime=%s",
				testconstants.SeveralVersionsDID,
				"2023-03-06T09:53:44Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceCollectionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceCollectionId=%s",
				testconstants.SeveralVersionsDID,
				"b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceVersion query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&resourceVersion=%s",
				testconstants.SeveralVersionsDID,
				"b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and checksum query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?metadata=true&checksum=%s",
				testconstants.SeveralVersionsDID,
				"b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceId=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"5e16a3f9-7c6e-4b6b-8e28-20f56780ee25",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceName query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceName=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"TestResource",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceType query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceType=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"TestType",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceVersionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceVersionTime=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"2023-03-06T09:53:44Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceMetadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceMetadata=true",
				testconstants.SeveralVersionsDID,
				"bar",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceCollectionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceCollectionId=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"2023-03-06T09:53:44Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceVersion query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&resourceVersion=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"1.0",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of service and checksum query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?service=%s&checksum=%s",
				testconstants.SeveralVersionsDID,
				"bar",
				"64ec88ca00b268e5ba1a35678a1b5316d212f4f366b2477232534a8aeca37f3c",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.SeveralVersionsDID,
						MethodSpecificId: "b5d70adf-31ca-4662-aa10-d3a54cd8f06c",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),
)
