//go:build integration

package metadata

import (
	"encoding/json"
	"fmt"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
	errors "github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Negative: Get DIDDoc with metadata query", func(testCase utils.NegativeTestCase) {
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
		"cannot get DIDDoc metadata with not supported metadata query parameter value",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=not_supported_value",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "representationNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionId query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true&versionId=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionTime query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true&versionTime=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				"2023-01-22T11:58:10.390039347Z",
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true&service=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentService,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "internalError",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.InternalErrorHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionId, versionTime, service query parameters",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true&versionId=%s&versionTime=%s&service=%s",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
				"2023-01-22T11:58:10.390039347Z",
				testconstants.NotExistentService,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: testconstants.UUIDStyleTestnetId,
						Method:           testconstants.ValidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),
)
