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

type resourceCollectionTestCase struct {
	didURL                      string
	resolutionType              types.ContentType
	expectedDereferencingResult *DereferencingResult
	expectedError               error
}

var _ = DescribeTable("Test ResourceCollectionEchoHandler function", func(testCase resourceCollectionTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDereferencingResult.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}

	expectedContentType := defineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

	err := resourceServices.ResourceCollectionEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error(), err.Error())
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
		"successful resolution",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", ValidDid),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				ContentStream: types.NewDereferencedResourceList(
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
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", NotExistDID),
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
		"invalid DID",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", InvalidDid),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        InvalidDid,
						MethodSpecificId: InvalidIdentifier,
						Method:           InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", ValidDid),
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
			expectedError: types.NewRepresentationNotSupportedError(ValidDid, types.JSON, nil, false),
		},
	),

	// TODO: add unit tests for
	// - redirect integration tests.
)
