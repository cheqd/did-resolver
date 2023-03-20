package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceDataTestCase struct {
	resolutionType   types.ContentType
	did              string
	resourceId       string
	expectedResource types.ContentStreamI
	expectedMetadata types.ResolutionDidDocMetadata
	expectedError    error
}

var validResourceDereferencing = types.DereferencedResourceData(validResource.Resource.Data)

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
	context, rec := setupContext(
		"/1.0/identifiers/:did/resources/:resource",
		[]string{"did", "resource"},
		[]string{testCase.did, testCase.resourceId},
		testCase.resolutionType,
		mockLedgerService)

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
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			resourceId:       ValidResourceId,
			expectedResource: &validResourceDereferencing,
			expectedMetadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		dereferenceResourceDataTestCase{
			resolutionType:   types.DIDJSONLD,
			did:              fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier),
			resourceId:       "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid representation",
		dereferenceResourceDataTestCase{
			resolutionType:   types.JSON,
			did:              ValidDid,
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
