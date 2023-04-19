//go:build unit

package request

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	didDocService "github.com/cheqd/did-resolver/services/diddoc"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Test Query handler with transformKey params", func(testCase queriesDIDDocTestCase) {
	request := httptest.NewRequest(http.MethodGet, testCase.didURL, nil)
	context, rec := utils.SetupEmptyContext(request, testCase.resolutionType, utils.MockLedger)
	expectedDIDResolution := testCase.expectedResolution.(*types.DidResolution)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		expectedDIDResolution.Did.Context = []string{types.DIDSchemaJSONLD, types.JsonWebKey2020JSONLD}
	} else if expectedDIDResolution.Did != nil {
		expectedDIDResolution.Did.Context = nil
	}

	expectedContentType := utils.DefineContentType(expectedDIDResolution.ResolutionMetadata.ContentType, testCase.resolutionType)

	err := didDocService.DidDocEchoHandler(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var resolutionResult types.DidResolution
		Expect(err).To(BeNil())
		Expect(json.Unmarshal(rec.Body.Bytes(), &resolutionResult)).To(BeNil())
		Expect(expectedDIDResolution.Did).To(Equal(resolutionResult.Did))
		Expect(expectedDIDResolution.Metadata).To(Equal(resolutionResult.Metadata))
		Expect(expectedContentType).To(Equal(resolutionResult.ResolutionMetadata.ContentType))
		Expect(expectedDIDResolution.ResolutionMetadata.DidProperties).To(Equal(resolutionResult.ResolutionMetadata.DidProperties))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2018 transformKey query parameter",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &types.DidDoc{
					Id: testconstants.ValidDIDDocResolution.Id,
					VerificationMethod: []types.VerificationMethod{
						{
							Id:              testconstants.ValidDIDDocResolution.VerificationMethod[0].Id,
							Type:            string(types.Ed25519VerificationKey2018),
							Controller:      testconstants.ValidDIDDocResolution.VerificationMethod[0].Controller,
							PublicKeyBase58: "6fYkiuzNvu5THPLV5PKc1b7NyCWQ9bJa2rnLhfRxiYUK",
						},
					},
					Service: testconstants.ValidDIDDocResolution.Service,
				},
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2020 transformKey query parameter",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2020,
			),
			resolutionType: types.DIDJSONLD,
			expectedResolution: &types.DidResolution{
				ResolutionMetadata: types.ResolutionMetadata{
					DidProperties: types.DidProperties{
						DidString:        testconstants.ValidDid,
						MethodSpecificId: testconstants.ValidIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				Did: &types.DidDoc{
					Id: testconstants.ValidDIDDocResolution.Id,
					VerificationMethod: []types.VerificationMethod{
						{
							Id:                 testconstants.ValidDIDDocResolution.VerificationMethod[0].Id,
							Type:               string(types.Ed25519VerificationKey2020),
							Controller:         testconstants.ValidDIDDocResolution.VerificationMethod[0].Controller,
							PublicKeyMultibase: "z6Mkk7ooKAEpGSZvPtBBkxHSrgfNnmnFZUYvishGXwPydmFh",
						},
					},
					Service: testconstants.ValidDIDDocResolution.Service,
				},
				Metadata: types.NewResolutionDidDocMetadata(
					testconstants.ValidDid, &testconstants.ValidMetadata,
					[]*resourceTypes.Metadata{testconstants.ValidResource.Metadata},
				),
			},
		},
	),

	Entry(
		"cannot get DIDDoc with not existent DID and supported transformKey query parameter",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s",
				testconstants.NotExistentTestnetDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewNotFoundError(
				testconstants.NotExistentTestnetDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc (JSONWebKey2020) with not supported transformKey query parameter",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=notSupportedTransformKey",
				testconstants.ValidDid,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and metadata query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&metadata=true",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and metadata query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&metadata=true",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceId query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceId=%s",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
				testconstants.ValidIdentifier,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceName query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceName=someName",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceType query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceType=someType",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceVersionTime query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceVersionTime=%s",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
				testconstants.NotEmptyTime.String(),
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceMetadata query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceMetadata=true",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceCollectionId query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceCollectionId=%s",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
				testconstants.ValidIdentifier,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"cannot get DIDDoc with combination of transformKey and resourceVersion query parameters",
		queriesDIDDocTestCase{
			didURL: fmt.Sprintf(
				"/1.0/identifiers/%s?transformKey=%s&resourceVersion=someVersion",
				testconstants.ValidDid,
				types.Ed25519VerificationKey2018,
			),
			resolutionType:     types.DIDJSONLD,
			expectedResolution: &types.DidResolution{},
			expectedError: types.NewRepresentationNotSupportedError(
				testconstants.ValidDid, types.DIDJSONLD, nil, false,
			),
		},
	),
)
