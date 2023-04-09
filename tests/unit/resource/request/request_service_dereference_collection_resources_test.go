//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
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
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDereferencingResult.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}

	expectedContentType := utils.DefineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

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
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.ValidDid),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: types.NewDereferencedResourceList(
					testconstants.ValidDid,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.NotExistentTestnetDid),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid DID",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.InvalidDID),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDID,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(testconstants.InvalidDID, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.ValidDid),
			resolutionType: types.JSON,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(testconstants.ValidDid, types.JSON, nil, false),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase utils.RedirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.DidURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.ResolutionType, utils.MockLedger)

	err := resourceServices.ResourceCollectionEchoHandler(context)
	if err != nil {
		Expect(testCase.ExpectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.ExpectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.ExpectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get collection of resources with an old 16 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL:                 fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.OldIndy16CharStyleTestnetDid),
			ResolutionType:         types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.MigratedIndy16CharStyleTestnetDid),
			ExpectedError:          nil,
		},
	),

	Entry(
		"can redirect when it try to get collection of resources with an old 32 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL:                 fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.OldIndy32CharStyleTestnetDid),
			ResolutionType:         types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.MigratedIndy32CharStyleTestnetDid),
			ExpectedError:          nil,
		},
	),
)
