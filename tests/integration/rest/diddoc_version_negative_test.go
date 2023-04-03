//go:build integration

package rest

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net/http"

// 	testconstants "github.com/cheqd/did-resolver/tests/constants"
// 	"github.com/cheqd/did-resolver/types"
// 	"github.com/go-resty/resty/v2"
// 	. "github.com/onsi/ginkgo/v2"
// 	. "github.com/onsi/gomega"
// )

// var _ = DescribeTable("Negative: Get DIDDoc version", func(testCase negativeTestCase) {
// 	client := resty.New()

// 	resp, err := client.R().
// 		SetHeader("Accept", testCase.resolutionType).
// 		Get(testCase.didURL)
// 	Expect(err).To(BeNil())

// 	var receivedDidDereferencing dereferencingResult
// 	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
// 	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))

// 	expectedDidDereferencing := testCase.expectedResult.(dereferencingResult)
// 	assertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
// },

// 	Entry(
// 		"cannot get DIDDoc with an not existent DID",
// 		negativeTestCase{
// 			didURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s/version/%s",
// 				testconstants.NotExistentMainnetDid,
// 				testconstants.ValidIdentifier,
// 			),
// 			expectedResult: dereferencingResult{
// 				Context: "",
// 				DereferencingMetadata: types.DereferencingMetadata{
// 					ContentType:     types.DIDJSONLD,
// 					ResolutionError: "notFound",
// 					DidProperties: types.DidProperties{
// 						DidString:        testconstants.NotExistentMainnetDid,
// 						MethodSpecificId: testconstants.NotExistentIdentifier,
// 						Method:           testconstants.ValidMethod,
// 					},
// 				},
// 				ContentStream: nil,
// 				Metadata:      types.ResolutionDidDocMetadata{},
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 	),

// 	Entry(
// 		"cannot get collection of resources with invalid DID",
// 		negativeTestCase{
// 			didURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s/version/%s",
// 				testconstants.InvalidDID,
// 				testconstants.ValidIdentifier,
// 			),
// 			resolutionType: testconstants.DefaultResolutionType,
// 			expectedResult: dereferencingResult{
// 				Context: "",
// 				DereferencingMetadata: types.DereferencingMetadata{
// 					ContentType:     types.DIDJSONLD,
// 					ResolutionError: "methodNotSupported",
// 					DidProperties: types.DidProperties{
// 						DidString:        testconstants.InvalidDID,
// 						MethodSpecificId: testconstants.InvalidIdentifier,
// 						Method:           testconstants.InvalidMethod,
// 					},
// 				},
// 				ContentStream: nil,
// 				Metadata:      types.ResolutionDidDocMetadata{},
// 			},
// 			expectedStatusCode: http.StatusNotImplemented,
// 		},
// 	),

// 	Entry(
// 		"cannot get DIDDoc with an existent DID, but not existent versionId",
// 		negativeTestCase{
// 			didURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s/version/%s",
// 				testconstants.IndyStyleMainnetDid,
// 				testconstants.NotExistentIdentifier,
// 			),
// 			resolutionType: testconstants.DefaultResolutionType,
// 			expectedResult: dereferencingResult{
// 				Context: "",
// 				DereferencingMetadata: types.DereferencingMetadata{
// 					ContentType:     types.DIDJSONLD,
// 					ResolutionError: "notFound",
// 					DidProperties: types.DidProperties{
// 						DidString:        testconstants.IndyStyleMainnetDid,
// 						MethodSpecificId: "Ps1ysXP2Ae6GBfxNhNQNKN",
// 						Method:           testconstants.ValidMethod,
// 					},
// 				},
// 				ContentStream: nil,
// 				Metadata:      types.ResolutionDidDocMetadata{},
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 	),

// 	Entry(
// 		"cannot get DIDDoc with an existent DID, but an invalid versionId",
// 		negativeTestCase{
// 			didURL: fmt.Sprintf(
// 				"http://localhost:8080/1.0/identifiers/%s/version/%s",
// 				testconstants.UUIDStyleMainnetDid,
// 				testconstants.InvalidIdentifier,
// 			),
// 			expectedResult: dereferencingResult{
// 				Context: "",
// 				DereferencingMetadata: types.DereferencingMetadata{
// 					ContentType:     types.DIDJSONLD,
// 					ResolutionError: "notFound",
// 					DidProperties: types.DidProperties{
// 						DidString:        testconstants.UUIDStyleMainnetDid,
// 						MethodSpecificId: "c82f2b02-bdab-4dd7-b833-3e143745d612",
// 						Method:           testconstants.ValidMethod,
// 					},
// 				},
// 				ContentStream: nil,
// 				Metadata:      types.ResolutionDidDocMetadata{},
// 			},
// 			expectedStatusCode: http.StatusNotFound,
// 		},
// 	),
// )
