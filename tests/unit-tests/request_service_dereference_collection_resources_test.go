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

type dereferenceCollectionResourcesTestCase struct {
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
		[]string{testCase.did},
		testCase.resolutionType,
		mockLedgerService)

	if (testCase.resolutionType == "" || testCase.resolutionType == types.DIDJSONLD) && testCase.expectedError == nil {
		testCase.expectedResource.AddContext(types.DIDSchemaJSONLD)
	} else if testCase.expectedResource != nil {
		testCase.expectedResource.RemoveContext()
	}

	expectedContentType := defineContentType(testCase.expectedResolutionType, testCase.resolutionType)

	err := resourceServices.ResourceCollectionEchoHandler(context)
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
			resolutionType:   types.DIDJSONLD,
			did:              fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier),
			expectedResource: nil,
			expectedMetadata: types.ResolutionDidDocMetadata{},
			expectedError: types.NewNotFoundError(
				fmt.Sprintf("did:%s:%s:%s", ValidMethod, ValidNamespace, NotExistIdentifier), types.DIDJSONLD, nil, false,
			),
		},
	),
)
