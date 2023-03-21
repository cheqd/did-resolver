package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceDataTestCase struct {
	didURL           string
	resolutionType   types.ContentType
	expectedResource types.ContentStreamI
	expectedError    error
}

var validResourceDereferencing = types.DereferencedResourceData(validResource.Resource.Data)

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
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
		dereferenceResourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, ValidResourceId),
			resolutionType:   types.DIDJSONLD,
			expectedResource: &validResourceDereferencing,
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		dereferenceResourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", NotExistDID, ValidResourceId),
			resolutionType:   types.DIDJSONLD,
			expectedResource: nil,
			expectedError:    types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		dereferenceResourceDataTestCase{
			didURL:           fmt.Sprintf("/1.0/identifiers/%s/resources/%s", ValidDid, ValidResourceId),
			resolutionType:   types.JSON,
			expectedResource: nil,
			expectedError:    types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
