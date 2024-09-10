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
	expectedDIDResolution *types.DidResolution
	expectedError         error
}

var _ = DescribeTable("Test DIDDocEchoHandler function", func(testCase resolveDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.LinkedDomainsJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDIDResolution.Did != nil {
		testCase.expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(testCase.expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

	err := didDocServices.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var resolutionResult types.DidResolution
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &resolutionResult)).To(BeNil())
		Expect(testCase.expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
		Expect(testCase.expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(testCase.expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"can get DIDDoc with an existent DID",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
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
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid DID method",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(
					testconstants.DIDStructure,
					testconstants.InvalidMethod,
					testconstants.ValidTestnetNamespace,
					testconstants.ValidIdentifier,
				),
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
				Metadata: types.ResolutionDidDocMetadata{},
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
				Metadata: types.ResolutionDidDocMetadata{},
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
				Metadata: types.ResolutionDidDocMetadata{},
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
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(testconstants.InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get DIDDoc with an invalid representation",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", testconstants.ExistentDid),
			resolutionType: types.JSON,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
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
