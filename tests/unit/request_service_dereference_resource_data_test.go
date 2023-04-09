package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceServices "github.com/cheqd/did-resolver/services/resource"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
)

type resourceDataTestCase struct {
	didURL           string
	resolutionType   types.ContentType
	expectedResource types.ContentStreamI
	expectedError    error
}

var validResourceDereferencing = types.DereferencedResourceData(validResource.Resource.Data)

var _ = DescribeTable("Test ResourceDataEchoHandler function", func(testCase resourceDataTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	expectedContentType := types.ContentType(validResource.Metadata.MediaType)

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
		"successful resolution",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, ValidResourceId),
			resolutionType:   types.DIDJSONLD,
			expectedResource: &validResourceDereferencing,
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", NotExistDID, ValidResourceId),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid DID",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", InvalidDid, ValidResourceId),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewMethodNotSupportedError(InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"a valid DID, but not existent resourceId",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, NotExistIdentifier),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"a valid DID, but an invalid resourceId",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, InvalidIdentifier),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewInvalidDIDUrlError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		resourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, ValidResourceId),
			resolutionType:   types.JSON,
			expectedResource: nil,
			expectedError:    types.NewRepresentationNotSupportedError(ValidDid, types.JSON, nil, true),
		},
	),
)

var _ = DescribeTable("Test redirect DID", func(testCase redirectDIDTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := setupEmptyContext(request, testCase.resolutionType, mockLedgerService)

	err := resourceServices.ResourceDataEchoHandler(context)
	if err != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(testCase.expectedError).To(BeNil())
		Expect(http.StatusMovedPermanently).To(Equal(rec.Code))
		Expect(testCase.expectedDidURLRedirect).To(Equal(rec.Header().Get(echo.HeaderLocation)))
	}
},

	Entry(
		"can redirect when it try to get resource data with an old 16 characters Indy style DID",
		redirectDIDTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType: types.DIDJSONLD,
			expectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.MigratedIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			expectedError: nil,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 32 characters Indy style DID",
		redirectDIDTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			resolutionType: types.DIDJSONLD,
			expectedDidURLRedirect: fmt.Sprintf(
				"/1.0/identifiers/%s/resources/%s",
				testconstants.MigratedIndy32CharStyleTestnetDid,
				"214b8b61-a861-416b-a7e4-45533af40ada",
			),
			expectedError: nil,
		},
	),
)
