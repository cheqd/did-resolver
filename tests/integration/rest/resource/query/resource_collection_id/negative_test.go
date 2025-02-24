//go:build integration

package resource_collection_id_test

import (
	"encoding/json"
	"fmt"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get Collection of Resources with collectionId query", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())

	expectedDidDereferencing := testCase.ExpectedResult.(utils.DereferencingResult)
	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},
	Entry(
		"cannot get collection of resources with only collectionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceCollectionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.UUIDStyleTestnetId,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode,
		},
	),
	Entry(
		"cannot get collection of resources with not existent collectionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceCollectionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get collection of resources with an invalid collectionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceCollectionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: utils.DereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      types.ResolutionDidDocMetadata{},
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode, // it should be invalidDidUrl
		},
	),
)
