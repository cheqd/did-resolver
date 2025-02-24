//go:build integration

package transformKeys

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

var _ = DescribeTable("Negative: Get DIDDoc with transformKeys query parameter", func(testCase utils.NegativeTestCase) {
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
		"cannot get DIDDoc with not existent DID and not supported transformKeys query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=EDDSA",
				testconstants.TestHostAddress,
				testconstants.NotExistentTestnetDid,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with not supported transformKeys query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=EDDSA",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKeys and metadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&metadata=true",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        didWithEd25519VerificationKey2018Key,
						MethodSpecificId: identifierDidWithEd25519VerificationKey2018Key,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: types.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKeys and resourceId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceId=%s",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceName query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceName=someName",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceType query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceType=someType",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceVersionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceVersionTime=2006-01-02",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceMetadata query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceMetadata=true",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceCollectionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceCollectionId=%s",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
		"cannot get DIDDoc with combination of transformKeys and resourceVersion query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?transformKeys=%s&resourceVersion=someVersion",
				testconstants.TestHostAddress,
				didWithEd25519VerificationKey2018Key,
				types.Ed25519VerificationKey2020,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
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
