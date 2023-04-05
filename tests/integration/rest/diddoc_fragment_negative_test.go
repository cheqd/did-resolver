//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get DID#fragment", func(testCase negativeTestCase) {
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
		"cannot get DIDDoc fragment with an existent DID, but not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.UUIDStyleMainnetDid+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc fragment with not existent DID and not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey-1",
				testconstants.NotExistentMainnetDid+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc fragment with not existent DID",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey1",
				testconstants.NotExistentTestnetDid+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc fragment with an invalid DID",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey1",
				testconstants.InvalidDID+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc fragment with existent DID, but not existent #fragment",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.IndyStyleTestnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.IndyStyleTestnetDid,
						MethodSpecificId: "73wnEyHhkhXiH1Nq7w5Kgq",
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
		"cannot get DIDDoc fragment with existent old 16 characters Indy style DID, but not existent #fragment",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.OldIndy16CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy16CharStyleTestnetDid,
						MethodSpecificId: "CpeMubv5yw63jXyrgRRsxR",
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
		"cannot get DIDDoc fragment with existent old 32 characters Indy style DID, but not existent #fragment",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.OldIndy32CharStyleTestnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy32CharStyleTestnetDid,
						MethodSpecificId: "3KpiDD6Hxs4i2G7FtpiGhu",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotFound,
		},
	),
)
