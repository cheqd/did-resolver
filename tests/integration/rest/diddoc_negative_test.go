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

type getDidDocNegativeTestCase struct {
	didURL                string
	resolutionType        string
	expectedDIDResolution *types.DidResolution
	expectedStatusCode    int
}

var _ = DescribeTable("Negative: Get DIDDoc", func(testCase getDidDocNegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDIDResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDIDResolution)).To(BeNil())

	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))
	Expect(testCase.expectedDIDResolution.Context).To(Equal(receivedDIDResolution.Context))
	Expect(testCase.expectedDIDResolution.ResolutionMetadata.ContentType).To(Equal(receivedDIDResolution.ResolutionMetadata.ContentType))
	Expect(testCase.expectedDIDResolution.ResolutionMetadata.ResolutionError).To(Equal(receivedDIDResolution.ResolutionMetadata.ResolutionError))
	Expect(testCase.expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(receivedDIDResolution.ResolutionMetadata.DidProperties))
	Expect(receivedDIDResolution.Did).To(BeZero())
	Expect(receivedDIDResolution.Metadata).To(BeZero())
},

	Entry(
		"cannot get DIDDoc with not existent mainnet DID",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.NotExistentMainnetDid),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentMainnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get DIDDoc with not existent testnet DID",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.NotExistentTestnetDid),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get DIDDoc with mainnet DID that contains an invalid method",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.MainnetDIDWithInvalidMethod),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MainnetDIDWithInvalidMethod,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
	),

	Entry(
		"cannot get DIDDoc with testnet DID that contains an invalid method",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.TestnetDIDWithInvalidMethod),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.TestnetDIDWithInvalidMethod,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
	),

	Entry(
		"cannot get DIDDoc with DID that contains an invalid namespace",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.DIDWithInvalidNamespace),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDid",
					DidProperties: types.DidProperties{
						DidString:        testconstants.DIDWithInvalidNamespace,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
	),

	Entry(
		"It cannot get DIDDoc with an invalid DID",
		getDidDocNegativeTestCase{
			didURL:         fmt.Sprintf("http://localhost:8080/1.0/identifiers/%s", testconstants.InvalidDID),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDIDResolution: &types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDID,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
	),
)
