package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

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
	did := fmt.Sprintf("did:%s:%s:%s", testCase.method, testCase.namespace, testCase.identifier)

	diddocService := services.NewDIDDocService("cheqd", testCase.ledgerService)
	expectedDIDProperties := types.DidProperties{
		DidString:        did,
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

	resolutionResult, err := diddocService.Resolve(did, "", testCase.resolutionType)
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
)
