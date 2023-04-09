package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
)

type DereferencingResult struct {
	DereferencingMetadata *types.DereferencingMetadata          `json:"dereferencingMetadata"`
	ContentStream         *types.DereferencedResourceListStruct `json:"contentStream"`
	Metadata              *types.ResolutionDidDocMetadata       `json:"contentMetadata"`
}

type resourceMetadataTestCase struct {
	didURL                      string
	resolutionType              types.ContentType
	expectedDereferencingResult *DereferencingResult
	expectedError               error
}

var _ = DescribeTable("Test ResourceMetadataEchoHandler function", func(testCase resourceMetadataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDereferencingResult.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}
	expectedContentType := defineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

	err := resourceServices.ResourceMetadataEchoHandler(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult DereferencingResult
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)

		Expect(err).To(BeNil())
		Expect(unmarshalErr).To(BeNil())
		Expect(testCase.expectedDereferencingResult.ContentStream, dereferencingResult.ContentStream)
		Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"successful resolution",
		resourceMetadataTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/resources/%s/metadata", ValidDid, ValidResourceId),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: types.NewDereferencedResourceListStruct(
					ValidDid,
					[]*resourceTypes.Metadata{validResource.Metadata},
				),
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resourceMetadataTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/resources/%s/metadata", NotExistDID, ValidResourceId),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        NotExistDID,
						MethodSpecificId: NotExistIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resourceMetadataTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/resources/%s/metadata", ValidDid, ValidResourceId),
			resolutionType: types.JSON,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),

	// TODO: add unit tests for invalid DID case.
)
