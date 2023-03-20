package tests

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	"github.com/cheqd/did-resolver/types"
)

type resolveDIDDocTestCase struct {
	resolutionType         types.ContentType
	did                    string
	expectedDID            *types.DidDoc
	expectedMetadata       types.ResolutionDidDocMetadata
	expectedResolutionType types.ContentType
	expectedError          error
}

var validDIDResolution = types.NewDidDoc(&validDIDDoc)

var _ = DescribeTable("Test ResolveDIDDoc method", func(testCase resolveDIDDocTestCase) {
	context, rec := setupContext(
		"/1.0/identifiers/:did",
		[]string{"did"}, []string{testCase.did},
		testCase.resolutionType,
		mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDID.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDID != nil {
		testCase.expectedDID.Context = nil
	}

	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := didDocServices.DidDocEchoHandler(context)
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
			resolutionType:   types.DIDJSONLD,
			did:              fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier),
			expectedDID:      nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),
)
