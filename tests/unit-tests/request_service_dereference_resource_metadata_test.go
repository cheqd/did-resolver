package tests

import (
	"encoding/json"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
)

type dereferenceResourceMetadataTestCase struct {
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
		[]string{testCase.did, testCase.resourceId},
		testCase.resolutionType,
		mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedResource.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedResource != nil {
		testCase.expectedResource.RemoveContext()
	}
	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := resourceServices.ResourceMetadataEchoHandler(context)

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
			resolutionType:   types.DIDJSONLD,
			did:              fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier),
			resourceId:       "a86f9cae-0902-4a7c-a144-96b60ced2fc9",
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),

	Entry(
		"invalid representation",
		dereferenceResourceMetadataTestCase{
			resolutionType:   types.JSON,
			did:              ValidDid,
			resourceId:       ValidResourceId,
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError:    types.NewRepresentationNotSupportedError(ValidDid, types.DIDJSONLD, nil, true),
		},
	),
)
