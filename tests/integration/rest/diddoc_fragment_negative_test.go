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

type getDidDocFragmentNegativeTestCase struct {
	didURL                      string
	resolutionType              string
	expectedDereferencingResult DereferencingResult
	expectedStatusCode          int
}

var _ = DescribeTable("Negative: Get DID#fragment", func(testCase getDidDocFragmentNegativeTestCase) {
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
		"cannot get DIDDoc with not existent mainnet DID#fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.NotExistentMainnetDid+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc with not existent testnet DID#fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%skey1",
				testconstants.NotExistentTestnetDid+url.PathEscape(testconstants.HashTag),
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
		"cannot get DIDDoc with existent 22 bytes INDY style mainnet DID, but not existent #fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.IndyStyleMainnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get DIDDoc with existent 22 bytes INDY style testnet DID, but not existent #fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.IndyStyleTestnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
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
		"cannot get DIDDoc with existent UUID style mainnet DID, but not existent #fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.UUIDStyleMainnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleMainnetDid,
						MethodSpecificId: "c82f2b02-bdab-4dd7-b833-3e143745d612",
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
		"cannot get DIDDoc with existent UUID style testnet DID, but not existent #fragment",
		getDidDocFragmentNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s%s",
				testconstants.UUIDStyleTestnetDid+url.PathEscape(testconstants.HashTag),
				testconstants.NotExistentFragment,
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
)
