//go:build unit

package request

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/types/known/timestamppb"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type QueriesDIDDocTestCase struct {
	didURL             string
	resolutionType     types.ContentType
	expectedResolution types.ResolutionResultI
	expectedError      error
}

type ResourceTestCase struct {
	didURL           string
	resolutionType   types.ContentType
	expectedResource types.ContentStreamI
	expectedError    error
}

var (
	Data                   = []byte("{\"attr\":[\"name\",\"age\"]}")
	Checksum               = sha256.New().Sum(Data)
	VersionId1             = uuid.New().String()
	VersionId2             = uuid.New().String()
	ResourceIdName1        = uuid.New().String()
	ResourceIdName2        = uuid.New().String()
	ResourceIdType1        = uuid.New().String()
	ResourceIdType2        = uuid.New().String()
	DidDocBeforeCreated, _ = time.Parse(time.RFC3339, "2021-08-23T08:59:00Z")
	DidDocCreated, _       = time.Parse(time.RFC3339, "2021-08-23T09:00:00Z")
	DidDocAfterCreated, _  = time.Parse(time.RFC3339, "2021-08-23T09:00:30Z")
	Resource1Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:01:00Z")
	Resource2Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:02:00Z")
	DidDocUpdated, _       = time.Parse(time.RFC3339, "2021-08-23T09:03:00Z")
	DidDocAfterUpdated, _  = time.Parse(time.RFC3339, "2021-08-23T09:03:30Z")
	Resource3Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:04:00Z")
	Resource4Created, _    = time.Parse(time.RFC3339, "2021-08-23T09:05:00Z")
)

var (
	ResourceName1   = generateResource(ResourceIdName1, "Name1", "string", "1", timestamppb.New(Resource1Created))
	ResourceName2   = generateResource(ResourceIdName2, "Name2", "string", "2", timestamppb.New(Resource2Created))
	ResourceType1   = generateResource(ResourceIdType1, "Name", "Type1", "3", timestamppb.New(Resource3Created))
	ResourceType2   = generateResource(ResourceIdType2, "Name", "Type2", "4", timestamppb.New(Resource4Created))
	DidDocMetadata1 = generateMetadata(
		VersionId1,
		timestamppb.New(DidDocCreated),
		nil,
	)
	DidDocMetadata2 = generateMetadata(
		VersionId2,
		timestamppb.New(DidDocCreated),
		timestamppb.New(DidDocUpdated),
	)
)

var MockLedger = utils.NewMockLedgerService(
	&testconstants.ValidDIDDoc,
	[]*didTypes.Metadata{
		&DidDocMetadata1,
		&DidDocMetadata2,
	},
	[]resourceTypes.ResourceWithMetadata{
		ResourceName1,
		ResourceName2,
		ResourceType1,
		ResourceType2,
	},
)

func generateResource(resourceId, name, rtype, version string, created *timestamppb.Timestamp) resourceTypes.ResourceWithMetadata {
	return resourceTypes.ResourceWithMetadata{
		Resource: &resourceTypes.Resource{
			Data: Data,
		},
		Metadata: &resourceTypes.Metadata{
			CollectionId: testconstants.ValidIdentifier,
			Id:           resourceId,
			Name:         name,
			ResourceType: rtype,
			MediaType:    "application/json",
			Checksum:     fmt.Sprintf("%x", Checksum),
			Created:      created,
			Version:      version,
		},
	}
}

func generateMetadata(versionId string, created, updated *timestamppb.Timestamp) didTypes.Metadata {
	return didTypes.Metadata{
		VersionId:   versionId,
		Deactivated: false,
		Created:     created,
		Updated:     updated,
	}
}

var _ = DescribeTable("Test Query handlers with versionId and versionTime params", func(testCase QueriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
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
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName1.Metadata},
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
						ResourceType1.Metadata,
						ResourceName2.Metadata,
						ResourceName1.Metadata},
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
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName1.Metadata},
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
						ResourceType1.Metadata,
						ResourceName2.Metadata,
						ResourceName1.Metadata},
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
					[]*resourceTypes.Metadata{ResourceName2.Metadata, ResourceName1.Metadata},
				),
			},
			expectedError: nil,
		},
	),

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
				Metadata: types.ResolutionDidDocMetadata{},
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
				Metadata: types.ResolutionDidDocMetadata{},
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
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(testconstants.InvalidVersionId, types.DIDJSONLD, nil, true),
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
				Metadata: types.ResolutionDidDocMetadata{},
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
			expectedError:      types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)

var _ = DescribeTable("Test Query handlers with service and relativeRef params", func(testCase QueriesDIDDocTestCase) {
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
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. relativeRef case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionId + Service case + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0]),
			expectedError:      nil,
		},
	),
	Entry(
		"Positive. VersionTime + Service case + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.ValidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: types.NewServiceResult(testconstants.ValidService.ServiceEndpoint[0] + "foo"),
			expectedError:      nil,
		},
	),

	// Negative Cases
	Entry(
		"Negative. Service not found",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo", testconstants.ValidDid, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionId=%s&service=%s", testconstants.ValidDid, testconstants.ValidVersionId, testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionId + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionId=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.ValidVersionId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?versionTime=%s&service=%s", testconstants.ValidDid, testconstants.CreatedAfter.Format(time.RFC3339), testconstants.InvalidServiceId),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Service not found + versionTime + relativeRef",
		QueriesDIDDocTestCase{
			didURL:             fmt.Sprintf("/1.0/identifiers/%s?service=%s&relativeRef=foo&versionTime=%s", testconstants.ValidDid, testconstants.InvalidServiceId, testconstants.CreatedAfter.Format(time.RFC3339)),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: nil,
			expectedError:      types.NewNotFoundError(testconstants.InvalidServiceId, types.DIDJSONLD, nil, true),
		},
	),
)

var _ = DescribeTable("Test Query handlers with resource params", func(testCase ResourceTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)
	expectedContentType := types.ContentType(testconstants.ValidResource[0].Metadata.MediaType)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedResource.GetBytes(), rec.Body.Bytes())
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	// Positive cases
	// Entry(
	// 	"Positive. ResourceId case. The first item",
	// 	ResourceTestCase{
	// 		didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceId=%s", testconstants.ValidDid, ResourceIdName1),
	// 		resolutionType:   types.DIDJSONLD,
	// 		expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
	// 		expectedError:    nil,
	// 	},
	// ),
	// Entry(
	// 	"Positive. ResourceId + CollectionId case. The first item",
	// 	ResourceTestCase{
	// 		didURL:           fmt.Sprintf("/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s", testconstants.ValidDid, ResourceIdName1, testconstants.ValidIdentifier),
	// 		resolutionType:   types.DIDJSONLD,
	// 		expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
	// 		expectedError:    nil,
	// 	},
	// ),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339)),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
	Entry(
		"Positive. ResourceId + CollectionId + ResourceName + ResourceType + ResourceVersion + ResourceVersionTime + Checksum case. The first item",
		ResourceTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceCollectionId=%s&resourceName=%s&resourceType=%s&resourceVersion=%s&resourceVersionTime=%s&checksum=%s",
				testconstants.ValidDid,
				ResourceIdName1,
				testconstants.ValidIdentifier,
				ResourceName1.Metadata.Name,
				ResourceName1.Metadata.ResourceType,
				ResourceName1.Metadata.Version,
				DidDocUpdated.Format(time.RFC3339),
				ResourceName1.Metadata.Checksum),
			resolutionType:   types.DIDJSONLD,
			expectedResource: types.NewDereferencedResourceData(ResourceName1.Resource.Data),
			expectedError:    nil,
		},
	),
)
