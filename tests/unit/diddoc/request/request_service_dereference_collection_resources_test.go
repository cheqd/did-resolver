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

	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type resourceCollectionTestCase struct {
	didURL                      string
	resolutionType              types.ContentType
	expectedDereferencingResult *DereferencingResult
	expectedError               *types.IdentityError
}

var _ = DescribeTable("Test DidDocResourceCollectionEchoHandler function", func(testCase resourceCollectionTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	expectedContentType := utils.DefineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

	err := didDocServices.DidDocResourceCollectionEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult DereferencingResult
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)).To(BeNil())
		Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"can get collection of resource with an existent DID",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.ExistentDid),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get collection of resources with not existent DID",
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
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get collection of resources with an invalid DID",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.InvalidDid),
			resolutionType: types.DIDJSONLD,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewMethodNotSupportedError(testconstants.InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resourceCollectionTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s/metadata", testconstants.ExistentDid),
			resolutionType: types.JSON,
			expectedDereferencingResult: &DereferencingResult{
				DereferencingMetadata: &types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Metadata: &types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewRepresentationNotSupportedError(testconstants.ExistentDid, types.JSON, nil, false),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase utils.RedirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.DidURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.ResolutionType, utils.MockLedger)

	err := didDocServices.DidDocResourceCollectionEchoHandler(context)
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
