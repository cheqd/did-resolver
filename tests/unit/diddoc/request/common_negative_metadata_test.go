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

var _ = DescribeTable("Negative tests for mixed cases", func(testCase ResourceMetadataTestCase) {
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
		"Negative. NotExistent versionId and resourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				uuid.New().String(),
				ResourceIdName1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. versionTime before DIDDocument was created and resourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				DidDocBeforeCreated.Format(time.RFC3339),
				ResourceIdName1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. NotExistentDid and versionTime and resourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.NotExistentTestnetDid,
				DidDocBeforeCreated.Format(time.RFC3339),
				ResourceIdName1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. wrong DID and versionTime and resourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.DidWithInvalidNamespace,
				DidDocBeforeCreated.Format(time.RFC3339),
				ResourceIdName1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewInvalidDidError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. NotExistent resourceId",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionId=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				uuid.New().String(),
				ResourceIdName1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
	Entry(
		"Negative. Valid versionTime but not existent resourceId for such time",
		ResourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?versionTime=%s&resourceId=%s&resourceMetadata=true",
				testconstants.ValidDid,
				DidDocAfterCreated.Format(time.RFC3339),
				ResourceIdType1,
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
				Metadata:      &types.DereferencedResource{},
			},
			expectedError: types.NewNotFoundError(testconstants.ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
