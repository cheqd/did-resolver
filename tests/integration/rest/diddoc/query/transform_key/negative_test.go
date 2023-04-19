//go:build integration

package transformKey

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

var IdentifierDIDWithEd25519VerificationKey2018Key = "d8ac0372-0d4b-413e-8ef5-8e8f07822b2c"

var _ = DescribeTable("", func(testCase utils.NegativeTestCase) {
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
		"cannot get DIDDoc with not supported transformKey query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and metadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&metadata=true",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceId=%s",
				DidWithEd25519VerificationKey2018Key,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceName query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceName=someName",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceType query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceType=someType",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceType query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceType=someType",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceVersionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceVersionTime=someVersionTime",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceMetadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceMetadata=true",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceCollectionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceCollectionId=%s",
				DidWithEd25519VerificationKey2018Key,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
		"cannot get DIDDoc with combination of transformKey and resourceVersion query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA&resourceVersion=someVersion",
				DidWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        DidWithEd25519VerificationKey2018Key,
						MethodSpecificId: IdentifierDIDWithEd25519VerificationKey2018Key,
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
