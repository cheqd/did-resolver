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

var identifierDidWithEd25519VerificationKey2018Key = "d8ac0372-0d4b-413e-8ef5-8e8f07822b2c"

var _ = DescribeTable("Negative: Get DIDDoc with transformKey query parameter", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())


	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))
	expectedDidResolution, ok := testCase.ExpectedResult.(types.DidResolution)
	if ok {
		var receivedDidResolution types.DidResolution
		Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
		utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
	} else {
		expectedDidDereferencing := testCase.ExpectedResult.(utils.DereferencingResult)
		var receivedDidDereferencing utils.DereferencingResult
		Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
		utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
	}
},

	Entry(
		"cannot get DIDDoc with not existent DID and not supported transformKey query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA",
				testconstants.NotExistentTestnetDid,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with not supported transformKey query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=EDDSA",
				didWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and metadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&metadata=true",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceId=%s",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceName=someName",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceType=someType",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceVersionTime=2006-01-02",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceMetadata=true",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceCollectionId=%s",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s&resourceVersion=someVersion",
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
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
