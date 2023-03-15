package tests

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type resolveDIDDocTestCase struct {
	ledgerService          MockLedgerService
	resolutionType         types.ContentType
	did                    string
	expectedDID            *types.DidDoc
	expectedMetadata       types.ResolutionDidDocMetadata
	expectedResolutionType types.ContentType
	expectedError          error
}

var validDIDResolution = types.NewDidDoc(&validDIDDoc)

var _ = DescribeTable("Test ResolveDIDDoc method", func(testCase resolveDIDDocTestCase) {
	context, rec := setupContext("/1.0/identifiers/:did", []string{"did"}, []string{testCase.did}, testCase.resolutionType)
	requestService := services.NewRequestService("cheqd", testCase.ledgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDID.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDID != nil {
		testCase.expectedDID.Context = nil
	}
	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := requestService.ResolveDIDDoc(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var resolutionResult types.DidResolution
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &resolutionResult)
		Expect(unmarshalErr).To(BeNil())
		Expect(err).To(BeNil())
		Expect(testCase.expectedDID).To(Equal(resolutionResult.Did))
		Expect(testCase.expectedMetadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"successful resolution",
		resolveDIDDocTestCase{
			ledgerService:    NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			expectedDID:      &validDIDResolution,
			expectedMetadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		resolveDIDDocTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)
