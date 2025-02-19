//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Test Query handlers with versionId and versionTime params", func(testCase QueriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	request.Header.Set("Accept", string(testCase.resolutionType))
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)
	expectedDIDResolution := testCase.expectedResolution.(*types.DidResolution)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if expectedDIDResolution.Did != nil {
		expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var resolutionResult types.DidResolution
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &resolutionResult)
		Expect(unmarshalErr).To(BeNil())
		Expect(err).To(BeNil())
		Expect(expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
		Expect(expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	// Negative cases
	Entry(
		"Negative. Wrong VersionId case",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionId=%s", testconstants.ValidDid, testconstants.InvalidVersionId),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. No VersionId case",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionId=%s", testconstants.ValidDid, uuid.New().String()),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. VersionTime wrong format",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, "TimeBeforeCreation"),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. VersionTime before created",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, DidDocBeforeCreated.Format(time.RFC3339)),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Query parameter is not supported",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?unsupportedQuery=%s", testconstants.ValidDid, "blabla"),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Unsupported Accept Header",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s", testconstants.ValidDid),
			resolutionType:     types.TEXT,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.JSON, nil, false),
		},
	),
	Entry(
		"Negative. Invalid value for metadata query",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?metadata=xxxx", testconstants.ValidDid),
			resolutionType:     types.JSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.JSONLD, nil, false),
		},
	),
	Entry(
		"Negative. Invalid value for resourceMetadata query",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?resourceMetadata=xxxx", testconstants.ValidDid),
			resolutionType:     types.JSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewInternalError(testconstants.ValidDid, types.JSONLD, nil, false),
		},
	),
	Entry(
		"Negative. Unsupported Accept Header",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s", testconstants.ValidDid),
			resolutionType:     types.TEXT,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.JSON, nil, false),
		},
	),
	Entry(
		"Negative. Invalid value for metadata query",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?metadata=xxxx", testconstants.ValidDid),
			resolutionType:     types.JSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.JSONLD, nil, false),
		},
	),
	Entry(
		"Negative. Invalid value for resourceMetadata query",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?resourceMetadata=xxxx", testconstants.ValidDid),
			resolutionType:     types.JSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewInternalError(testconstants.ValidDid, types.JSONLD, nil, false),
		},
	),
)
