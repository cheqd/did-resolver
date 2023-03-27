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

type resourceMetadataNegativeTestCase struct {
	didURL                      string
	resolutionType              string
	expectedDereferencingResult DereferencingResult
	expectedStatusCode          int
}

var _ = DescribeTable("Negative: get resource metadata", func(testCase resourceMetadataNegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())

	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))
	Expect(testCase.expectedDereferencingResult.Context).To(Equal(receivedDidDereferencing.Context))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType).To(Equal(receivedDidDereferencing.DereferencingMetadata.ContentType))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.ResolutionError).To(Equal(receivedDidDereferencing.DereferencingMetadata.ResolutionError))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(receivedDidDereferencing.DereferencingMetadata.DidProperties))
	Expect(testCase.expectedDereferencingResult.ContentStream).To(Equal(receivedDidDereferencing.ContentStream))
	Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(receivedDidDereferencing.Metadata))
},

	Entry(
		"cannot get resource metadata with not existent DID and a valid resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get resource metadata with an existent DID and not existent resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: "c1685ca0-1f5b-439c-8eb8-5c0e85ab7cd0",
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
		"cannot get resource metadata with not existent DID and resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.NotExistentTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get resource metadata with an invalid DID and a valid resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.InvalidDID,
				testconstants.ValidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get resource metadata with an existent DID and an invalid resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.IndyStyleMainnetDid,
				testconstants.InvalidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get resource metadata with an invalid DID and resourceId",
		resourceMetadataNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.InvalidDID,
				testconstants.InvalidIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
)
