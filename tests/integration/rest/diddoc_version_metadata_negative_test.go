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

var _ = DescribeTable("Negative: Get DIDDoc version metadata", func(testCase negativeTestCase) {
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
		"cannot get DIDDoc version metadata with an existent DID and versionId, but not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
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
		"cannot get DIDDoc version metadata with not existent DID and not supported resolutionType",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
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
		"cannot get DIDDoc version metadata with not existent DID",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
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
		"cannot get DIDDoc version metadata with an invalid DID",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
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
		"cannot get DIDDoc version metadata with an existent DID, but not existent versionId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.IndyStyleMainnetDid,
				testconstants.NotExistentIdentifier,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedResult: dereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.IndyStyleMainnetDid,
						MethodSpecificId: "Ps1ysXP2Ae6GBfxNhNQNKN",
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
		"cannot get DIDDoc version metadata with an existent old 16 characters Indy style DID, but not existent versionId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
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
		"cannot get DIDDoc version metadata with an existent old 32 characters Indy style DID, but not existent versionId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
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

	Entry(
		"cannot get DIDDoc version metadata with an existent DID, but an invalid versionId",
		negativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.UUIDStyleMainnetDid,
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
