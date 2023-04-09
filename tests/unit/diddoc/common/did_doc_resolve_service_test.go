//go:build unit

package common

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

type resolveTestCase struct {
	did                   string
	resolutionType        types.ContentType
	expectedDIDResolution *types.DidResolution
	expectedError         *types.IdentityError
}

var _ = DescribeTable("Test Resolve method", func(testCase resolveTestCase) {
	diddocService := services.NewDIDDocService("cheqd", utils.MockLedger)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDIDResolution.Did != nil {
		testCase.expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(
		testCase.expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType,
	)

	resolutionResult, err := diddocService.Resolve(testCase.did, "", testCase.resolutionType)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Code).To(Equal(err.Code))
		Expect(testCase.expectedError.Message).To(Equal(err.Message))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
		Expect(testCase.expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(testCase.expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
	}
},

	Entry(
		"Successful resolution",
		resolveTestCase{
			did:            testconstants.ValidDid,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &testconstants.ValidDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resolveTestCase{
			did:            testconstants.NotExistentTestnetDid,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentTestnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false),
		},
	),
)
