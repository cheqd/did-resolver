package tests

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceDataTestCase struct {
	ledgerService    MockLedgerService
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
		testCase.ledgerService)

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
			ledgerService:    NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
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
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			resourceId:       "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid representation",
		dereferenceResourceDataTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.JSON,
			did:              ValidDid,
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
