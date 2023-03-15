package tests

import (
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type dereferencingTestCase struct {
	ledgerService         MockLedgerService
	dereferencingType     types.ContentType
	did                   string
	fragmentId            string
	queries               url.Values
	expectedContentStream types.ContentStreamI
	expectedMetadata      types.ResolutionDidDocMetadata
	expectedContentType   types.ContentType
	expectedError         *types.IdentityError
}

var _ = DescribeTable("Test Dereferencing method", func(testCase dereferencingTestCase) {
	diddocService := services.NewDIDDocService("cheqd", testCase.ledgerService)
	var expectedDIDProperties types.DidProperties
	if testCase.expectedError == nil {
		expectedDIDProperties = types.DidProperties{
			DidString:        ValidDid,
			MethodSpecificId: ValidIdentifier,
			Method:           ValidMethod,
		}
	}

	expectedContentType := testCase.expectedContentType
	if expectedContentType == "" {
		expectedContentType = testCase.dereferencingType
	}

	result, err := diddocService.ProcessDIDRequest(testCase.did, testCase.fragmentId, testCase.queries, nil, testCase.dereferencingType)
	dereferencingResult, _ := result.(*types.DidDereferencing)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedContentStream).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedMetadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))

		Expect(dereferencingResult.DereferencingMetadata.ResolutionError).To(BeEmpty())
		Expect(expectedDIDProperties).To(Equal(dereferencingResult.DereferencingMetadata.DidProperties))
	}
},

	Entry(
		"successful Secondary dereferencing (key)",
		dereferencingTestCase{
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			did:                   ValidDid,
			fragmentId:            validVerificationMethod.Id,
			expectedContentStream: types.NewVerificationMethod(&validVerificationMethod),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
	),

	Entry(
		"successful Secondary dereferencing (service)",
		dereferencingTestCase{
			ledgerService:         NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			dereferencingType:     types.DIDJSON,
			did:                   ValidDid,
			fragmentId:            validService.Id,
			expectedContentStream: types.NewService(&validService),
			expectedMetadata:      validFragmentMetadata,
			expectedError:         nil,
		},
	),

	Entry(
		"not supported query",
		dereferencingTestCase{
			ledgerService:         NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   ValidDid,
			queries:               validQuery,
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"key not found",
		dereferencingTestCase{
			ledgerService:         NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			dereferencingType:     types.DIDJSONLD,
			did:                   ValidDid,
			fragmentId:            "notFoundKey",
			expectedContentStream: nil,
			expectedMetadata:      types.ResolutionDidDocMetadata{},
			expectedError:         types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)
