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
	resolutionType        types.ContentType
	did                   string
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

	expectedContentType := testCase.expectedDIDResolution.ResolutionMetadata.ContentType
	if expectedContentType == "" {
		expectedContentType = testCase.resolutionType
	}

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
			resolutionType: types.DIDJSONLD,
			did:            ValidDid,
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        ValidDid,
						MethodSpecificId: ValidIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      &validDIDDocResolution,
				Metadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			},
			expectedError: nil,
		},
	),

	Entry(
		"DID not found",
		resolveTestCase{
			resolutionType: types.DIDJSONLD,
			did:            fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, NotExistIdentifier),
			expectedDIDResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, NotExistIdentifier),
						MethodSpecificId: NotExistIdentifier,
						Method:           ValidMethod,
					},
				},
				Did:      nil,
				Metadata: types.ResolutionDidDocMetadata{},
			},
			expectedError: types.NewNotFoundError(fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, NotExistIdentifier), types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid DID",
		resolveTestCase{
			resolutionType: types.DIDJSONLD,
			did:            InvalidDid,
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
			resolutionType: types.DIDJSONLD,
			did:            "did:" + InvalidMethod + ":" + ValidNamespace + ":" + ValidIdentifier,
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
			expectedError: types.NewNotFoundError(fmt.Sprintf("did:%s:%s:%s", InvalidMethod, ValidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid namespace",
		resolveTestCase{
			resolutionType: types.DIDJSONLD,
			did:            "did:" + ValidMethod + ":" + InvalidNamespace + ":" + ValidIdentifier,
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
			expectedError: types.NewNotFoundError(fmt.Sprintf("did:%s:%s:%s", ValidMethod, InvalidNamespace, ValidIdentifier), types.DIDJSONLD, nil, false),
		},
	),

	Entry(
		"invalid identifier",
		resolveTestCase{
			resolutionType: types.DIDJSONLD,
			did:            "did:" + ValidMethod + ":" + ValidNamespace + ":" + InvalidIdentifier,
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
			expectedError: types.NewNotFoundError(fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, InvalidIdentifier), types.DIDJSONLD, nil, false),
		},
	),
)
