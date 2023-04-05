//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get resource data", func(testCase negativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing dereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

	expectedDidDereferencing := testCase.expectedResult.(dereferencingResult)
	assertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"cannot get resource data with an existent DID, but not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.UUIDStyleMainnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType: string(types.JSON),
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotAcceptable,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType: string(types.JSON),
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotAcceptable,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and a valid resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
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
			expectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get resource data with an invalid DID and not existent resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.InvalidDID,
				testconstants.ValidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
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
			expectedStatusCode: http.StatusNotImplemented,
		},
	),

	Entry(
		"cannot get resource data with not existent DID and resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.NotExistentTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
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
			expectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get resource data with an existent DID and an invalid resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.IndyStyleMainnetDid,
				testconstants.InvalidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	),

	Entry(
		"cannot get resource data with an existent old 16 characters Indy style DID and an invalid resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	),

	Entry(
		"cannot get resource data with an existent old 32 characters Indy style DID and an invalid resourceId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	),
)
