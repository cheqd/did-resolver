//go:build integration

package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type resourceCollectionNegativeTestCase struct {
	didURL                      string
	resolutionType              string
	expectedDereferencingResult DereferencingResult
	expectedStatusCode          int
}

var _ = DescribeTable("Positive: Get collection of resources", func(testCase resourceCollectionNegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())

	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))
	Expect(testCase.expectedDereferencingResult.Context).To(Equal(receivedDidDereferencing.Context))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.ContentType).To(Equal(receivedDidDereferencing.DereferencingMetadata.ContentType))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.ResolutionError).To(Equal(receivedDidDereferencing.DereferencingMetadata.ResolutionError))
	Expect(testCase.expectedDereferencingResult.DereferencingMetadata.DidProperties).To(Equal(receivedDidDereferencing.DereferencingMetadata.DidProperties))
	Expect(testCase.expectedDereferencingResult.ContentStream).To(Equal(receivedDidDereferencing.ContentStream))
	Expect(testCase.expectedDereferencingResult.Metadata).To(Equal(receivedDidDereferencing.Metadata))
},

	Entry(
		"cannot get collection of resources with not existent DID",
		resourceCollectionNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.NotExistentMainnetDid,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentMainnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotFound,
		},
	),

	Entry(
		"cannot get collection of resources with invalid DID",
		resourceCollectionNegativeTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.InvalidDID,
			),
			resolutionType: testconstants.DefaultResolutionType,
			expectedDereferencingResult: DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDID,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			expectedStatusCode: http.StatusNotImplemented,
		},
	),
)
