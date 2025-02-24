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

	var receivedDidDereferencing types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())

	expectedDidDereferencing := testCase.ExpectedResult.(types.DidResolution)
	utils.AssertDidResolution(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"cannot get collection of resources with not existent collectionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?collectionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode,
		},
	),

	Entry(
		"cannot get collection of resources with an invalid collectionId query parameter",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?collectionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: string(types.DIDJSONLD),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.DIDJSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: types.InvalidDidUrlHttpCode, // it should be invalidDidUrl
		},
	),
)
