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

var _ = DescribeTable("Test resource negative cases with Metadata field", func(testCase ResourceMetadataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, MockLedger)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDereferencingResult.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}
	expectedContentType := utils.DefineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

	err := didDocService.DidDocEchoHandler(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult DereferencingResult
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)).To(BeNil())
		Expect(testCase.expectedDereferencingResult.ContentStream).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"Negative. ResourceId not found",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				uuid.New().String(),
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceId wrong format",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				"SomeNotUUID",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. resourceVersionTime wrong format",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.ValidDid,
				"SomeNotUUID",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceCollectionId not found",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceCollectionId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				uuid.New().String(),
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceType is not found",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceType=%s&resourceMetadata=true",
				testconstants.ValidDid,
				"NotExistentType",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceName is not found",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceName=%s&resourceMetadata=true",
				testconstants.ValidDid,
				"NotExistentName",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceVersion is not found",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersion=%s&resourceMetadata=true",
				testconstants.ValidDid,
				"NotExistentVersion",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceVersion + ResourceMetadata=false",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersion=%s&resourceMetadata=false",
				testconstants.ValidDid,
				ResourceType1.Metadata.Version,
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. checksum wrong",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceId=%s&checksum=%s&resourceMetadata=true",
				testconstants.ValidDid,
				ResourceIdName1,
				"wrongChecksum",
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. ResourceVersionTime before the first resource created",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=true",
				testconstants.ValidDid,
				DidDocBeforeCreated.Format(time.RFC3339),
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Megative. ResourceVersionTime and resourceMetadata=false",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceVersionTime=%s&resourceMetadata=false",
				testconstants.ValidDid,
				Resource2Created.Format(time.RFC3339),
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. resourceCollectionId and resourceMetadata=false",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceCollectionId=%s&resourceMetadata=false",
				testconstants.ValidDid,
				ResourceName1.Metadata.CollectionId,
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. resourceType and resourceMetadata=false",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?resourceType=%s&resourceMetadata=false",
				testconstants.ValidDid,
				ResourceType12.Metadata.ResourceType,
			),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: &types.DereferencedResourceListStruct{},
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
