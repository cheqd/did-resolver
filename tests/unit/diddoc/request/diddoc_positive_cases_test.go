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

var _ = DescribeTable("Test Query handlers with versionId and versionTime params", func(testCase QueriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)
	expectedDIDResolution := testCase.expectedResolution.(*types.DidResolution)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.LinkedDomainsJSONLD, types.JsonWebKey2020JSONLD}
	} else if expectedDIDResolution.Did != nil {
		expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

	err := didDocService.DidDocEchoHandler(context)
	var resolutionResult types.DidResolution
	unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &resolutionResult)
	Expect(unmarshalErr).To(BeNil())
	Expect(err).To(BeNil())
	Expect(expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
	Expect(expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
	Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
	Expect(expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
	Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
},

	// Positive cases
	Entry(
		"Positive. VersionId case. The first item",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionId=%s", testconstants.ValidDid, VersionId1),
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
					testconstants.ValidDid, &DidDocMetadata1,
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName12.Metadata, ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId case. The second item. All the resources should be returned",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionId=%s", testconstants.ValidDid, VersionId2),
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
					testconstants.ValidDid, &DidDocMetadata2,
					[]*resourceTypes.Metadata{
						ResourceType2.Metadata,
						ResourceType12.Metadata,
						ResourceType1.Metadata,
						ResourceName2.Metadata,
						ResourceName12.Metadata,
						ResourceName1.Metadata,
					},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionTime case. Case after creation but before update",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, DidDocAfterCreated.Format(time.RFC3339Nano)),
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
					testconstants.ValidDid, &DidDocMetadata1,
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName12.Metadata, ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionTime case. Case after updating",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s", testconstants.ValidDid, DidDocAfterUpdated.Format(time.RFC3339Nano)),
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
					testconstants.ValidDid, &DidDocMetadata2,
					[]*resourceTypes.Metadata{
						ResourceType2.Metadata,
						ResourceType12.Metadata,
						ResourceType1.Metadata,
						ResourceName2.Metadata,
						ResourceName12.Metadata,
						ResourceName1.Metadata,
					},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionTime and VersionId case. Case after creation but before update",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&versionId=%s", testconstants.ValidDid, DidDocAfterCreated.Format(time.RFC3339Nano), VersionId1),
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
					testconstants.ValidDid, &DidDocMetadata1,
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName12.Metadata, ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. Metadata = false. Should return all the resources in metadata",
		QueriesDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s?metadata=false", testconstants.ValidDid),
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
					testconstants.ValidDid, &DidDocMetadata2,
					[]*resourceTypes.Metadata{
						ResourceType2.Metadata,
						ResourceType12.Metadata,
						ResourceType1.Metadata,
						ResourceName2.Metadata,
						ResourceName12.Metadata,
						ResourceName1.Metadata,
					},
				),
			},
			expectedError: nil,
		},
	),
)
