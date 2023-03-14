package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
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

type dereferenceResourceDataTestCase struct {
	ledgerService    MockLedgerService
	resolutionType   types.ContentType
	did              string
	resourceId       string
	expectedResource types.ContentStreamI
	expectedMetadata types.ResolutionDidDocMetadata
	expectedError    error
}

var validResourceDereferencing = types.DereferencedResourceData(validResource.Resource.Data)

var _ = DescribeTable("Test DereferenceResourceData method", func(testCase dereferenceResourceDataTestCase) {
	context, rec := setupContext(
		"/1.0/identifiers/:did/resources/:resource",
		[]string{"did", "resource"},
		[]string{testCase.did, testCase.resourceId}, testCase.resolutionType)
	requestService := services.NewRequestService("cheqd", testCase.ledgerService)
	expectedContentType := types.ContentType(validResource.Metadata.MediaType)

	err := requestService.DereferenceResourceData(context)
	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		Expect(err).To(BeNil())
		Expect(testCase.expectedResource.GetBytes(), rec.Body.Bytes())
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"successful resolution",
		dereferenceResourceDataTestCase{
			ledgerService:    NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			resourceId:       ValidResourceId,
			expectedResource: &validResourceDereferencing,
			expectedMetadata: types.NewResolutionDidDocMetadata(ValidDid, &validMetadata, []*resourceTypes.Metadata{validResource.Metadata}),
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		dereferenceResourceDataTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			resourceId:       "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)

type dereferenceResourceMetadataTestCase struct {
	ledgerService          MockLedgerService
	resolutionType         types.ContentType
	did                    string
	resourceId             string
	expectedResource       *types.DereferencedResourceList
	expectedMetadata       types.ResolutionDidDocMetadata
	expectedResolutionType types.ContentType
	expectedError          error
}

var _ = DescribeTable("Test DereferenceResourceMetadata method", func(testCase dereferenceResourceMetadataTestCase) {
	context, rec := setupContext(
		"/1.0/identifiers/:did/resources/:resource/metadata",
		[]string{"did", "resource"},
		[]string{testCase.did, testCase.resourceId}, testCase.resolutionType)
	requestService := services.NewRequestService("cheqd", testCase.ledgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedResource.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedResource != nil {
		testCase.expectedResource.RemoveContext()
	}
	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := requestService.DereferenceResourceMetadata(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error()).To(Equal(err.Error()))
	} else {
		var dereferencingResult struct {
			DereferencingMetadata types.DereferencingMetadata    `json:"dereferencingMetadata"`
			ContentStream         types.DereferencedResourceList `json:"contentStream"`
			Metadata              types.ResolutionDidDocMetadata `json:"contentMetadata"`
		}
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)

		Expect(err).To(BeNil())
		Expect(unmarshalErr).To(BeNil())
		Expect(*testCase.expectedResource, dereferencingResult.ContentStream)
		Expect(testCase.expectedMetadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},
	Entry(
		"successful resolution",
		dereferenceResourceMetadataTestCase{
			ledgerService:  NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType: types.DIDJSONLD,
			did:            ValidDid,
			resourceId:     ValidResourceId,
			expectedResource: types.NewDereferencedResourceList(
				ValidDid,
				[]*resourceTypes.Metadata{validResource.Metadata},
			),
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		dereferenceResourceMetadataTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			resourceId:       "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)

type dereferenceCollectionResourcesTestCase struct {
	ledgerService          MockLedgerService
	resolutionType         types.ContentType
	did                    string
	expectedResource       *types.DereferencedResourceList
	expectedMetadata       types.ResolutionDidDocMetadata
	expectedResolutionType types.ContentType
	expectedError          error
}

var _ = DescribeTable("Test DereferenceCollectionResources method", func(testCase dereferenceCollectionResourcesTestCase) {
	context, rec := setupContext(
		"/1.0/identifiers/:did/metadata",
		[]string{"did"},
		[]string{testCase.did}, testCase.resolutionType)
	requestService := services.NewRequestService("cheqd", testCase.ledgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedResource.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedResource != nil {
		testCase.expectedResource.RemoveContext()
	}
	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := requestService.DereferenceCollectionResources(context)

	if testCase.expectedError != nil {
		Expect(testCase.expectedError.Error(), err.Error())
	} else {
		var dereferencingResult struct {
			DereferencingMetadata types.DereferencingMetadata    `json:"dereferencingMetadata"`
			ContentStream         types.DereferencedResourceList `json:"contentStream"`
			Metadata              types.ResolutionDidDocMetadata `json:"contentMetadata"`
		}
		unmarshalErr := json.Unmarshal(rec.Body.Bytes(), &dereferencingResult)

		Expect(err).To(BeNil())
		Expect(unmarshalErr).To(BeNil())
		Expect(*testCase.expectedResource).To(Equal(dereferencingResult.ContentStream))
		Expect(testCase.expectedMetadata).To(Equal(dereferencingResult.Metadata))
		Expect(expectedContentType).To(Equal(dereferencingResult.DereferencingMetadata.ContentType))
		Expect(expectedContentType).To(Equal(types.ContentType(rec.Header().Get("Content-Type"))))
	}
},

	Entry(
		"successful resolution",
		dereferenceCollectionResourcesTestCase{
			ledgerService:  NewMockLedgerService(&validDIDDoc, &validMetadata, &validResource),
			resolutionType: types.DIDJSONLD,
			did:            ValidDid,
			expectedResource: types.NewDereferencedResourceList(
				ValidDid,
				[]*resourceTypes.Metadata{validResource.Metadata},
			),
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    nil,
		},
	),

	Entry(
		"DID not found",
		dereferenceCollectionResourcesTestCase{
			ledgerService:    NewMockLedgerService(&didTypes.DidDoc{}, &didTypes.Metadata{}, &resourceTypes.ResourceWithMetadata{}),
			resolutionType:   types.DIDJSONLD,
			did:              ValidDid,
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewNotFoundError(ValidDid, types.DIDJSONLD, nil, false),
		},
	),
)

func defineContentType(expectedContentType types.ContentType, resolutionType types.ContentType) types.ContentType {
	if expectedContentType == "" {
		return resolutionType
	}
	return expectedContentType
}

func setupContext(path string, paramsNames []string, paramsValues []string, resolutionType types.ContentType) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	context := e.NewContext(req, rec)
	context.SetPath(path)
	context.SetParamNames(paramsNames...)
	context.SetParamValues(paramsValues...)
	req.Header.Add("accept", string(resolutionType))
	return context, rec
}
