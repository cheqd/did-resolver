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

	resourceServices "github.com/cheqd/did-resolver/services/resource"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
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
	expectedDereferencingResult *types.ResourceDereferencing
	expectedError               error
}

var _ = DescribeTable("Test ResourceMetadataEchoHandler function", func(testCase resourceMetadataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	if testCase.expectedDereferencingResult.ContentStream != nil {
		testCase.expectedDereferencingResult.ContentStream.RemoveContext()
	}
	expectedContentType := utils.DefineContentType(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType, testCase.resolutionType)

	err := resourceServices.ResourceMetadataEchoHandler(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult *types.ResourceDereferencing
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)).To(BeNil())
		Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"can get resource metadata with an existent DID and resourceId",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.ExistentDid,
				testconstants.ExistentResourceId,
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
					ContentMetadata: types.NewDereferencedResource(
						testconstants.ExistentDid,
						testconstants.ValidResource[0].Metadata,
					),
				},
			},
			expectedError: nil,
		},
	),

	Entry(
		"cannot get resource metadata with not existent DID",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.NotExistentTestnetDid,
				testconstants.ExistentResourceId,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      nil,
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource metadata with an invalid DID",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.InvalidDid,
				testconstants.ExistentResourceId,
			),
			resolutionType: types.JSONLD,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      nil,
			},
			expectedError: types.NewMethodNotSupportedError(testconstants.InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource metadata with an existent DID, but not existent resourceId",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.ExistentDid,
				testconstants.NotExistentIdentifier,
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
				Metadata:      nil,
			},
			expectedError: types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"cannot get resource metadata with an existent DID, but an invalid resourceId",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.ExistentDid,
				testconstants.InvalidIdentifier,
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
				Metadata:      nil,
			},
			expectedError: types.NewInvalidDidUrlError(testconstants.ExistentDid, types.DIDJSONLD, nil, true),
		},
	),

	Entry(
		"cannot get resource metadata with an invalid content type",
		resourceMetadataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.ExistentDid,
				testconstants.ExistentResourceId,
			),
			resolutionType: types.TEXT,
			expectedDereferencingResult: &types.ResourceDereferencing{
				DereferencingMetadata: types.DereferencingMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ExistentDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      nil,
			},
			expectedError: types.NewRepresentationNotSupportedError(testconstants.ExistentDid, types.JSON, nil, false),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase utils.RedirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.DidURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.ResolutionType, utils.MockLedger)

	err := resourceServices.ResourceMetadataEchoHandler(context)
	if err != nil {
		Expect(testCase.ExpectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.ExpectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.ExpectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get resource metadata with an old 16 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.MigratedIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ExpectedError: nil,
		},
	),

	Entry(
		"can redirect when it try to get resource metadata with an old 32 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			ResolutionType: types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.MigratedIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			ExpectedError: nil,
		},
	),
)
