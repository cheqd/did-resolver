//go:build integration

package metadata_test

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

var _ = DescribeTable("Negative: get resource metadata", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedResourceDereferencing utils.ResourceDereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedResourceDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	expectedResourceDereferencing := testCase.ExpectedResult.(utils.ResourceDereferencingResult)
	utils.AssertResourceMetadata(expectedResourceDereferencing, receivedResourceDereferencing)
},

	Entry(
		"cannot get resource metadata with an existent DID and resourceId, but not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.JSON),
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with not existent DID and not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.JSON),
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSON,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with not existent DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentMainnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with an invalid DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.InvalidDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.MethodNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with an existent DID and not existent resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.UUIDStyleTestnetDid,
						MethodSpecificId: "c1685ca0-1f5b-439c-8eb8-5c0e85ab7cd0",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with an existent old 16 characters Indy style DID and not existent resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy16CharStyleTestnetDid,
						MethodSpecificId: "CpeMubv5yw63jXyrgRRsxR",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with an existent old 32 characters Indy style DID and not existent resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy32CharStyleTestnetDid,
						MethodSpecificId: "3KpiDD6Hxs4i2G7FtpiGhu",
						Method:           testconstants.ValidMethod,
					},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.NotFoundHttpCode,
		},
	),

	Entry(
		"cannot get resource metadata with an existent DID and an invalid resourceId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: utils.ResourceDereferencingResult{
				Context: "",
				DereferencingMetadata: types.DereferencingMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				ContentStream: nil,
				Metadata:      &types.DereferencedResource{},
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),
)
