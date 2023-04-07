//go:build integration

package data

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
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
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusNotAcceptable,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusNotAcceptable,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and a valid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get resource data with an invalid DID and not existent resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.InvalidDID,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDID,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: http.StatusNotImplemented,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get resource data with an existent DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusBadRequest,
		},
	),

	Entry(
		"cannot get resource data with an existent old 16 characters Indy style DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusBadRequest,
		},
	),

	Entry(
		"cannot get resource data with an existent old 32 characters Indy style DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
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
			ExpectedStatusCode: http.StatusBadRequest,
		},
	),
)