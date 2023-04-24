//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type queriesDIDDocTestCase struct {
	didURL             string
	resolutionType     types.ContentType
	expectedResolution types.ResolutionResultI
	expectedError      error
}

var _ = DescribeTable("Test Query handlers with versionId and versionTime params", func(testCase queriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)
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

	// Positive cases
	Entry(
		"Positive. VersionId case",
		queriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionId=%s", testconstants.ValidDid, testconstants.ValidVersionId),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionTime case",
		queriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339)),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	// Negative cases
	Entry(
		"Negative VersionId case",
		queriesDIDDocTestCase{
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
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. VersionTime before created",
		queriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, testconstants.CreatedBefore.Format(time.RFC3339)),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
			expectedError: types.NewNotFoundError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Query parameter is not supported",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?unsupportedQuery=%s", testconstants.ValidDid, "blabla"),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError:      types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)

var _ = DescribeTable("Test Query handlers with service and relativeRef params", func(testCase queriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(rec.Code).To(Equal(http.StatusSeeOther))
		Expect(string(testCase.expectedResolution.GetBytes())).To(Equal(context.Response().Header().Get("Location")))
		Expect(err).To(BeNil())
	}
},

	// Positive cases
	Entry(
		"Positive. Service case",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. relativeRef case",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case + relativeRef",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case + relativeRef",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),

	// Negative Cases
	Entry(
		"Negative. Service not found",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + relativeRef",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId + relativeRef",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionId=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.ValidVersionId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime + relativeRef",
		queriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionTime=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.CreatedAfter.Format(time.RFC3339)),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
)
