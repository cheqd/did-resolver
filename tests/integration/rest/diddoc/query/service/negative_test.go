//go:build integration

package service

import (
	. "github.com/onsi/ginkgo/v2"
)

var SeveralVersionsDIDIdentifier = "b5d70adf-31ca-4662-aa10-d3a54cd8f06c"

var _ = DescribeTable("Negative: Get DIDDoc with service query", func() {

},

	Entry(
		"cannot redirect to serviceEndpoint with not existent service query parameter",
		// utils.NegativeTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?service=%s",
		// 		SeveralVersionsDID,
		// 		testconstants.NotExistentService,
		// 	),
		// 	ResolutionType: testconstants.DefaultResolutionType,
		// 	ExpectedResult: types.DidResolution{
		// 		Context: "",
		// 		ResolutionMetadata: types.ResolutionMetadata{
		// 			ContentType:     types.DIDJSONLD,
		// 			ResolutionError: "notFound",
		// 			DidProperties: types.DidProperties{
		// 				DidString:        SeveralVersionsDID,
		// 				MethodSpecificId: SeveralVersionsDIDIdentifier,
		// 				Method:           testconstants.ValidMethod,
		// 			},
		// 		},
		// 		Did:      nil,
		// 		Metadata: types.ResolutionDidDocMetadata{},
		// 	},
		// 	ExpectedStatusCode: errors.NotFoundHttpCode,
		// },
	),

	Entry(
		"cannot redirect to serviceEndpoint with relativeRef query parameter",
		// utils.NegativeTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?relativeRef=/info",
		// 		SeveralVersionsDID,
		// 	),
		// 	ResolutionType: testconstants.DefaultResolutionType,
		// 	ExpectedResult: types.DidResolution{
		// 		Context: "",
		// 		ResolutionMetadata: types.ResolutionMetadata{
		// 			ContentType:     types.DIDJSONLD,
		// 			ResolutionError: "representationNotSupported",
		// 			DidProperties: types.DidProperties{
		// 				DidString:        SeveralVersionsDID,
		// 				MethodSpecificId: SeveralVersionsDIDIdentifier,
		// 				Method:           testconstants.ValidMethod,
		// 			},
		// 		},
		// 		Did:      nil,
		// 		Metadata: types.ResolutionDidDocMetadata{},
		// 	},
		// 	ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		// },
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but an invalid relativeRef URI parameters",
	// 	utils.NegativeTestCase{
	// 		DidURL: fmt.Sprintf(
	// 			"http://localhost:8080/1.0/identifiers/%s?relativeRef=/info",
	// 			SeveralVersionsDID,
	// 		),
	// 		ResolutionType: testconstants.DefaultResolutionType,
	// 		ExpectedResult: types.DidResolution{
	// 			Context: "",
	// 			ResolutionMetadata: types.ResolutionMetadata{
	// 				ContentType:     types.DIDJSONLD,
	// 				ResolutionError: "representationNotSupported",
	// 				DidProperties: types.DidProperties{
	// 					DidString:        SeveralVersionsDID,
	// 					MethodSpecificId: SeveralVersionsDIDIdentifier,
	// 					Method:           testconstants.ValidMethod,
	// 				},
	// 			},
	// 			Did:      nil,
	// 			Metadata: types.ResolutionDidDocMetadata{},
	// 		},
	// 		ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
	// 	},
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionId query parameters",
		// utils.NegativeTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?service=%s&versionId=%s",
		// 		SeveralVersionsDID,
		// 		"bar",
		// 		testconstants.NotExistentIdentifier,
		// 	),
		// 	ResolutionType: testconstants.DefaultResolutionType,
		// 	ExpectedResult: types.DidResolution{
		// 		Context: "",
		// 		ResolutionMetadata: types.ResolutionMetadata{
		// 			ContentType:     types.DIDJSONLD,
		// 			ResolutionError: "notFound",
		// 			DidProperties: types.DidProperties{
		// 				DidString:        SeveralVersionsDID,
		// 				MethodSpecificId: SeveralVersionsDIDIdentifier,
		// 				Method:           testconstants.ValidMethod,
		// 			},
		// 		},
		// 		Did:      nil,
		// 		Metadata: types.ResolutionDidDocMetadata{},
		// 	},
		// 	ExpectedStatusCode: errors.NotFoundHttpCode,
		// },
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionTime query parameters",
		// utils.NegativeTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?service=%s&versionTime=%s",
		// 		SeveralVersionsDID,
		// 		"bar",
		// 		"2023-03-06T09:36:56.56204903Z",
		// 	),
		// 	ResolutionType: testconstants.DefaultResolutionType,
		// 	ExpectedResult: types.DidResolution{
		// 		Context: "",
		// 		ResolutionMetadata: types.ResolutionMetadata{
		// 			ContentType:     types.DIDJSONLD,
		// 			ResolutionError: "notFound",
		// 			DidProperties: types.DidProperties{
		// 				DidString:        SeveralVersionsDID,
		// 				MethodSpecificId: SeveralVersionsDIDIdentifier,
		// 				Method:           testconstants.ValidMethod,
		// 			},
		// 		},
		// 		Did:      nil,
		// 		Metadata: types.ResolutionDidDocMetadata{},
		// 	},
		// 	ExpectedStatusCode: errors.NotFoundHttpCode,
		// },
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionId and versionTime query parameters",
		// utils.NegativeTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?service=%s&versionId=%s&versionTime=%s",
		// 		SeveralVersionsDID,
		// 		"bar",
		// 		testconstants.NotExistentIdentifier,
		// 		"2023-03-06T09:36:56.56204903Z",
		// 	),
		// 	ResolutionType: testconstants.DefaultResolutionType,
		// 	ExpectedResult: types.DidResolution{
		// 		Context: "",
		// 		ResolutionMetadata: types.ResolutionMetadata{
		// 			ContentType:     types.DIDJSONLD,
		// 			ResolutionError: "notFound",
		// 			DidProperties: types.DidProperties{
		// 				DidString:        SeveralVersionsDID,
		// 				MethodSpecificId: SeveralVersionsDIDIdentifier,
		// 				Method:           testconstants.ValidMethod,
		// 			},
		// 		},
		// 		Did:      nil,
		// 		Metadata: types.ResolutionDidDocMetadata{},
		// 	},
		// 	ExpectedStatusCode: errors.NotFoundHttpCode,
		// },
	),
)