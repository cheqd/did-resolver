//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type resolveDIDDocTestCase struct {
	didURL                string
	resolutionType        types.ContentType
	acceptHeader          string
	expectedDIDResolution interface{} // interface to accept type DidResolution as well as DidDocument
	expectedError         error
}

var _ = DescribeTable("Test DIDDocEchoHandler function", func(testCase resolveDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	request.Header.Set("Accept", testCase.acceptHeader) // Set Accept header dynamically
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	var didDoc *types.DidDoc
	var didResolution *types.DidResolution

	switch v := testCase.expectedDIDResolution.(type) {
	case *types.DidDoc:
		didDoc = v
	case *types.DidResolution:
		didResolution = v
	}
	if didResolution != nil {
		if (testCase.resolutionType == "" || testCase.resolutionType == types.JSONLD) && testCase.expectedError == nil {
			didResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.LinkedDomainsJSONLD, types.JsonWebKey2020JSONLD}
		}
	} else if didDoc != nil {
		didDoc.Context = nil
	}

	expectedContentType := utils.DefineContentType(
		func() types.ContentType {
			if didResolution != nil {
				return didResolution.ResolutionMetadata.ContentType
			}
			return ""
		}(),
		testCase.resolutionType,
	)
	responseContentType := utils.ResponseContentType(testCase.acceptHeader, false)

	err := didDocServices.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(err).To(BeNil())
		if didDoc != nil {
			var resolutionResult types.DidDoc
			Expect(json.Unmarshal(rec.Body.Bytes(), &resolutionResult)).To(BeNil())
			Expect(didDoc).To(Equal(resolutionResult))
			Expect(expectedContentType).To(Equal(rec.Header().Get("Content-Type")))
		} else if didResolution != nil {
			var resolutionResult types.DidResolution
			Expect(json.Unmarshal(rec.Body.Bytes(), &resolutionResult)).To(BeNil())
			Expect(didResolution.Did).To(Equal(resolutionResult.Did))
			Expect(didResolution.Metadata).To(Equal(resolutionResult.Metadata))
			Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
			Expect(didResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
			Expect(responseContentType).To(Equal(rec.Header().Get("Content-Type")))
		}
	}
},
	Entry(
		"can get DIDDoc with an existent DID",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType: types.JSONLD,
			acceptHeader:   string(types.JSONLD) + ";profile=\"" + types.W3IDDIDRES + "\"",
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType: types.DIDRES,
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ExistentDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource[0].Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"can get DIDDoc with an existent DID and application/did accept header",
		resolveDIDDocTestCase{
			didURL:                fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType:        types.DIDRES,
			acceptHeader:          string(types.DIDRES),
			expectedDIDResolution: testconstants.ValidDIDDocResolution,
			expectedError:         nil,
		},
	),
	Entry(
		"can get DIDDoc with an existent DID and application/did+json accept header",
		resolveDIDDocTestCase{
			didURL:                fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType:        types.DIDJSON,
			acceptHeader:          string(types.DIDJSON),
			expectedDIDResolution: testconstants.ValidDIDDocResolution,
			expectedError:         nil,
		},
	),
	Entry(
		"can get DIDDoc with an existent DID and application/did+ld+json accept header",
		resolveDIDDocTestCase{
			didURL:                fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType:        types.DIDJSONLD,
			acceptHeader:          string(types.DIDJSONLD),
			expectedDIDResolution: testconstants.ValidDIDDocResolution,
			expectedError:         nil,
		},
	),

	Entry(
		"cannot get DIDDoc with not existent DID",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.NotExistentTestnetDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid DID method",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				testconstants.TestnetDidWithInvalidMethod,
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString: fmt.Sprintf(
							testconstants.DIDStructure,
							testconstants.InvalidMethod,
							testconstants.ValidTestnetNamespace,
							testconstants.ValidIdentifier,
						),
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.InvalidMethod,
					testconstants.ValidTestnetNamespace,
					testconstants.ValidIdentifier,
				),
				types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid DID namespace",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.ValidMethod,
					testconstants.InvalidNamespace,
					testconstants.ValidIdentifier,
				),
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString: fmt.Sprintf(
							testconstants.DIDStructure,
							testconstants.ValidMethod,
							testconstants.InvalidNamespace,
							testconstants.ValidIdentifier,
						),
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidError(
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.ValidMethod,
					testconstants.InvalidNamespace,
					testconstants.ValidIdentifier,
				),
				types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid DID identifier",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.ValidMethod,
					testconstants.ValidTestnetNamespace,
					testconstants.InvalidIdentifier,
				),
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString: fmt.Sprintf(
							testconstants.DIDStructure,
							testconstants.ValidMethod,
							testconstants.ValidTestnetNamespace,
							testconstants.InvalidIdentifier,
						),
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidError(
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.ValidMethod,
					testconstants.ValidTestnetNamespace,
					testconstants.InvalidIdentifier,
				),
				types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid DID",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.InvalidDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(testconstants.InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid representation",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType: types.TEXT,
			acceptHeader:   string(types.TEXT),
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(testconstants.ExistentDid, types.JSON, nil, false),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase utils.RedirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.DidURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.ResolutionType, utils.MockLedger)

	err := didDocServices.DidDocEchoHandler(context)
	if err != nil {
		Expect(testCase.ExpectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.ExpectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.ExpectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get DIDDoc with an old 16 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL:                 fmt.Sprintf("/1.0/identifiers/%s", testconstants.OldIndy16CharStyleTestnetDid),
			ResolutionType:         types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s", testconstants.MigratedIndy16CharStyleTestnetDid),
			ExpectedError:          nil,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc with an old 32 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL:                 fmt.Sprintf("/1.0/identifiers/%s", testconstants.OldIndy32CharStyleTestnetDid),
			ResolutionType:         types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s", testconstants.MigratedIndy32CharStyleTestnetDid),
			ExpectedError:          nil,
		},
	),
)
