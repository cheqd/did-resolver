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
	didURL                string
	resolutionType        types.ContentType
	expectedDIDResolution *types.DidResolution
	expectedError         *types.IdentityError
}

var _ = DescribeTable("Test DIDDocEchoHandler method", func(testCase resolveDIDDocTestCase) {
	context, rec := setupContext(
		testCase.didURL,
		[]string{"did"},
		[]string{getDID(testCase.didURL)},
		testCase.resolutionType,
		mockLedgerService,
	)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDIDResolution.Did != nil {
		testCase.expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := defineContentType(testCase.expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

	err := didDocServices.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var resolutionResult types.DidResolution
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &resolutionResult)
		Expect(unmarshalErr).To(BeNil())
		Expect(err).To(BeNil())
		Expect(testCase.expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
		Expect(testCase.expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"successful resolution",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", ValidDid),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did: &validDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					ValidDid, &validMetadata,
					[]*resourceTypes.Metadata{validResource.Metadata},
				),
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resolveDIDDocTestCase{
			didURL:         fmt.Sprintf("/1.0/identifiers/%s", NotExistDID),
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        NotExistDID,
						MethodSpecificId: NotExistIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(NotExistDID, types.DIDJSONLD, nil, false),
		},
	),
)
