//go:build unit

package request

import (
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

type resourceDataTestCase struct {
	didURL           string
	resolutionType   types.ContentType
	expectedResource types.ContentStreamI
	expectedError    error
}

var _ = DescribeTable("Test ResourceDataEchoHandler function", func(testCase resourceDataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)

	expectedContentType := types.ContentType(testconstants.ValidResource[0].Metadata.MediaType)

	err := resourceServices.ResourceDataEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedResource.GetBytes(), rec.Body.Bytes())
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"can get resource data with an existent DID and resourceId",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.ExistentDid,
				testconstants.ExistentResourceId,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: &testconstants.ValidResourceDereferencing,
			expectedError:    nil,
		},
	),

	Entry(
		"cannot get resource data with not existent DID",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.NotExistentTestnetDid,
				testconstants.ExistentResourceId,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource data with an invalid DID",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.InvalidDid,
				testconstants.ExistentResourceId,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewMethodNotSupportedError(testconstants.InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource data with an existent DID, but not existent resourceId",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.ExistentDid,
				testconstants.NotExistentIdentifier,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(testconstants.ExistentDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource data with an existent DID, but an invalid resourceId",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.ExistentDid,
				testconstants.InvalidIdentifier,
			),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDidUrlError(testconstants.ExistentDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"cannot get resource data with invalid representation",
		resourceDataTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.ExistentDid,
				testconstants.ExistentResourceId,
			),
			resolutionType:   types.TEXT,
			expectedResource: &testconstants.ValidResourceDereferencing,
			expectedError:    nil,
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase utils.RedirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.DidURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.ResolutionType, utils.MockLedger)

	err := resourceServices.ResourceDataEchoHandler(context)
	if err != nil {
		Expect(testCase.ExpectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.ExpectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.ExpectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get resource data with an old 16 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.MigratedIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ExpectedError: nil,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 32 characters Indy style DID",
		utils.RedirectDIDTestCase{
			DidURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			ResolutionType: types.DIDJSONLD,
			ExpectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.MigratedIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			ExpectedError: nil,
		},
	),
)
