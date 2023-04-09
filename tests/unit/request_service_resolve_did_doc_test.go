package tests

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
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDIDResolution.Did != nil {
		testCase.expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := defineContentType(testCase.expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

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
		"successful resolution",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", ValidDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did: &validDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					ValidDid, &validMetadata,
					[]*resourceTypes.Metadata{validResource.Metadata},
				),
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", NotExistDID),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        NotExistDID,
						MethodSpecificId: NotExistIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid DID method",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(testconstants.DIDStructure, InvalidMethod, ValidNamespace, ValidIdentifier),
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        fmt.Sprintf(testconstants.DIDStructure, InvalidMethod, ValidNamespace, ValidIdentifier),
						MethodSpecificId: ValidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(
				fmt.Sprintf(testconstants.DIDStructure, InvalidMethod, ValidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid DID namespace",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(testconstants.DIDStructure, ValidMethod, InvalidNamespace, ValidIdentifier),
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        fmt.Sprintf(testconstants.DIDStructure, ValidMethod, InvalidNamespace, ValidIdentifier),
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDIDError(
				fmt.Sprintf(testconstants.DIDStructure, ValidMethod, InvalidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid DID identifier",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s",
				fmt.Sprintf(testconstants.DIDStructure, ValidMethod, ValidNamespace, InvalidIdentifier),
			),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        fmt.Sprintf(testconstants.DIDStructure, ValidMethod, ValidNamespace, InvalidIdentifier),
						MethodSpecificId: InvalidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDIDError(
				fmt.Sprintf(testconstants.DIDStructure, ValidMethod, ValidNamespace, InvalidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid DID",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", InvalidDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        InvalidDid,
						MethodSpecificId: InvalidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", ValidDid),
			resolutionType: types.JSON,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(ValidDid, types.JSON, nil, false),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase redirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	err := didDocServices.DidDocEchoHandler(context)
	if err != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.expectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get DIDDoc with an old 16 characters Indy style DID",
		redirectDIDTestCase{
			didURL:                 fmt.Sprintf("/1.0/identifiers/%s", testconstants.OldIndy16CharStyleTestnetDid),
			resolutionType:         types.DIDJSONLD,
			expectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s", testconstants.MigratedIndy16CharStyleTestnetDid),
			expectedError:          nil,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc with an old 32 characters Indy style DID",
		redirectDIDTestCase{
			didURL:                 fmt.Sprintf("/1.0/identifiers/%s", testconstants.OldIndy32CharStyleTestnetDid),
			resolutionType:         types.DIDJSONLD,
			expectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s", testconstants.MigratedIndy32CharStyleTestnetDid),
			expectedError:          nil,
		},
	),
)
