//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Tests for mixed DidDoc and resource cases", func(testCase resolveDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	request.Header.Set("Accept", testCase.acceptHeader) // Set Accept header dynamically
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)
	expectedDIDResolution := testCase.expectedDIDResolution

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.LinkedDomainsJSONLD, types.JsonWebKey2020JSONLD}
	} else if expectedDIDResolution.Did != nil {
		expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)
	responseContentType := utils.ResponseContentType(testCase.acceptHeader, false)

	err := didDocService.DidDocEchoHandler(context)
	var resolutionResult types.DidResolution
	unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &resolutionResult)
	Expect(unmarshalErr).To(BeNil())
	Expect(err).To(BeNil())
	testCase.expectedDIDResolution.Did.Context = resolutionResult.Did.Context
	Expect(expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
	Expect(expectedDIDResolution.Metadata.Resources).To(Equal(resolutionResult.Metadata.Resources))
	Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
	Expect(expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
	Expect(responseContentType).To(Equal(rec.Header().Get("Content-Type")))
},
	Entry(
		"Positive. VersionId + VersionTime + ResourceId",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceCollectionId get all resources",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceCollectionId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceName1.Metadata.CollectionId,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{
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
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType: types.DIDRES,
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339),
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType: types.DIDRES,
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime + Checksum",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s&checksum=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339),
				ResourceName1.Metadata.Checksum,
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType: types.DIDRES,
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceVersionTime return resources something between",
		resolveDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceVersionTime=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				Resource2Created.Format(time.RFC3339),
			),
			acceptHeader:   string(types.JSONLD) + ";profile=" + types.W3IDDIDRES,
			resolutionType: types.JSONLD,
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
					testconstants.ValidDid,
					&testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{
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
