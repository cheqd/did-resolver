package tests

import (
	"fmt"
	"net/url"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

var _ = Describe("Test GetDIDFragment method", func() {
	DIDDoc := types.NewDidDoc(&validDIDDoc)

	It("can find a existent fragment in VerificationMethod", func() {
		fragmentId := DIDDoc.VerificationMethod[0].Id
		expectedFragment := &DIDDoc.VerificationMethod[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, DIDDoc)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("can find a existent fragment in Service", func() {
		fragmentId := DIDDoc.Service[0].Id
		expectedFragment := &DIDDoc.Service[0]

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, DIDDoc)
		Expect(fragment).To(Equal(expectedFragment))
	})

	It("cannot find a not-existent fragment", func() {
		fragmentId := "fake_id"

		didDocService := services.DIDDocService{}

		fragment := didDocService.GetDIDFragment(fragmentId, DIDDoc)
		Expect(fragment).To(BeNil())
	})
})

type resolveTestCase struct {
	ledgerService          MockLedgerService
	resolutionType         types.ContentType
	identifier             string
	method                 string
	namespace              string
	expectedDID            *types.DidDoc
	expectedMetadata       types.ResolutionDidDocMetadata
	expectedResolutionType types.ContentType
	expectedError          *types.IdentityError
}

var _ = DescribeTable("Test Resolve method", func(testCase resolveTestCase) {
	id := fmt.Sprintf("did:%s:%s:%s", testCase.method, testCase.namespace, testCase.identifier)

	diddocService := services.NewDIDDocService("cheqd", testCase.ledgerService)
	expectedDIDProperties := types.DidProperties{
		DidString:        id,
		MethodSpecificId: testCase.identifier,
		Method:           testCase.method,
	}

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDID.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDID != nil {
		testCase.expectedDID.Context = nil
	}

	expectedContentType := testCase.expectedResolutionType
	if expectedContentType == "" {
		expectedContentType = testCase.resolutionType
	}

	resolutionResult, err := diddocService.Resolve(id, "", testCase.resolutionType)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedDID).To(Equal(resolutionResult.Did))
		Expect(testCase.expectedMetadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(expectedDIDProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
	}
},

	Entry(
		"Successful resolution",
		resolveTestCase{
			ledgerService:    NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      &validDIDDocResolution,
			expectedMetadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		resolveTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid DID",
		resolveTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       "oooooo0000OOOO_invalid_did",
			method:           ValidMethod,
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid method",
		resolveTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           "not_supported_method",
			namespace:        ValidNamespace,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid namespace",
		resolveTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			identifier:       ValidIdentifier,
			method:           ValidMethod,
			namespace:        "invalid_namespace",
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"representation is not supported",
		resolveTestCase{
			ledgerService:          NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:         "text/html,application/xhtml+xml",
			identifier:             ValidIdentifier,
			method:                 ValidMethod,
			namespace:              ValidNamespace,
			expectedDID:            nil,
			expectedMetadata:       types.ResolutionDidDocMetadata{},
			expectedResolutionType: types.JSON,
			expectedError:          types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
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
