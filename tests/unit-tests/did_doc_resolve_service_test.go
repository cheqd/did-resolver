package tests

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type resolveTestCase struct {
	did                   string
	resolutionType        types.ContentType
	expectedDIDResolution *types.DidResolution
	expectedError         *types.IdentityError
}

var _ = DescribeTable("Test Resolve method", func(testCase resolveTestCase) {
	diddocService := services.NewDIDDocService("cheqd", mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if testCase.expectedDIDResolution.Did != nil {
		testCase.expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := defineContentType(
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
			did:            ValidDid,
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
		resolveTestCase{
			did:            NotExistDID,
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

	Entry(
		"invalid DID",
		resolveTestCase{
			did:            InvalidDid,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        InvalidDid,
						MethodSpecificId: InvalidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(InvalidDid, types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid method",
		resolveTestCase{
			did:            "did:" + InvalidMethod + ":" + ValidNamespace + ":" + ValidIdentifier,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        "did:" + InvalidMethod + ":" + ValidNamespace + ":" + ValidIdentifier,
						MethodSpecificId: ValidIdentifier,
						Method:           InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid namespace",
		resolveTestCase{
			did:            "did:" + ValidMethod + ":" + InvalidNamespace + ":" + ValidIdentifier,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        "did:" + ValidMethod + ":" + InvalidNamespace + ":" + ValidIdentifier,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid identifier",
		resolveTestCase{
			did:            "did:" + ValidMethod + ":" + ValidNamespace + ":" + InvalidIdentifier,
			resolutionType: types.DIDJSONLD,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        "did:" + ValidMethod + ":" + ValidNamespace + ":" + InvalidIdentifier,
						MethodSpecificId: InvalidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),
)
