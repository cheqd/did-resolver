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

var _ = DescribeTable("Tests for mixed DidDoc and resource cases", func(testCase ResourceMetadataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)

	if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}
	expectedContentType := utils.DefineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)
	responseContentType := utils.ResponseContentType(request.Header.Get("accept"), true)

	err := didDocService.DidDocEchoHandler(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult *types.ResourceDereferencing
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)).To(BeNil())
		Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(responseContentType).To(Equal(rec.Header().Get("Content-Type")))
	}
},
	Entry(
		"Positive. VersionId + VersionTime + ResourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceCollectionId get all resources",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceCollectionId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceName1.Metadata.CollectionId,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{
							ResourceName2.Metadata,
							ResourceName12.Metadata,
							ResourceName1.Metadata,
						},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				ResourceIdName1,
				ResourceName1.Metadata.CollectionId,
				ResourceName1.Metadata.Name,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType",
		ResourceMetadataTestCase{
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
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion",
		ResourceMetadataTestCase{
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
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime",
		ResourceMetadataTestCase{
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
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceId + ResourceCollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime + Checksum",
		ResourceMetadataTestCase{
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
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{ResourceName1.Metadata},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
	Entry(
		"Positive. VersionId + VersionTime + ResourceVersionTime return resources something between",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&versionTime=%s&resourceVersionTime=%s&resourceMetadata=true",
				testconstants.ValidDid,
				VersionId1,
				DidDocUpdated.Format(time.RFC3339Nano),
				Resource2Created.Format(time.RFC3339),
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata: &types.ResolutionResourceMetadata{
					Resources: &types.NewDereferencedResourceListStruct(
						testconstants.ValidDid,
						[]*resourceTypes.Metadata{
							ResourceName2.Metadata,
							ResourceName12.Metadata,
							ResourceName1.Metadata,
						},
					).Resources,
				},
			},
			expectedError: nil,
		},
	),
)
