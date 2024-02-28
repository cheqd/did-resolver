//go:build integration

package data_test

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

var _ = DescribeTable("Negative: Get resource data", func(testCase utils.NegativeTestCase) {
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
		"cannot get resource data with an existent DID, but not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.JSON),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.JSON),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and a valid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentMainnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource data with an invalid DID and not existent resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.InvalidDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.MethodNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.NotExistentTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource data with an existent DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get resource data with an existent old 16 characters Indy style DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get resource data with an existent old 32 characters Indy style DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),
)
