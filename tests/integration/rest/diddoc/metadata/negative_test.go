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

var _ = DescribeTable("Negative: Get DIDDoc version metadata", func(testCase utils.NegativeTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	expectedDidResolution := testCase.ExpectedResult.(types.DidResolution)
	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},

	Entry(
		"cannot get DIDDoc version metadata with an existent DID and versionId, but not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.TEXT),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.TEXT,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc version metadata with not existent DID and not supported ResolutionType",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: string(types.TEXT),
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.TEXT,
					ResolutionError: "representationNotSupported",
					DidProperties:   types.DidProperties{},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.RepresentationNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc version metadata with not existent DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.NotExistentMainnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.NotExistentMainnetDid,
						MethodSpecificId: testconstants.NotExistentIdentifier,
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
		"cannot get DIDDoc version metadata with an invalid DID",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.InvalidDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "methodNotSupported",
					DidProperties: types.DidProperties{
						DidString:        testconstants.InvalidDid,
						MethodSpecificId: testconstants.InvalidIdentifier,
						Method:           testconstants.InvalidMethod,
					},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.MethodNotSupportedHttpCode,
		},
	),

	Entry(
		"cannot get DIDDoc version metadata with an existent DID, but not existent versionId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.IndyStyleMainnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.IndyStyleMainnetDid,
						MethodSpecificId: "Ps1ysXP2Ae6GBfxNhNQNKN",
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
		"cannot get DIDDoc version metadata with an existent old 16 characters Indy style DID, but not existent versionId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy16CharStyleTestnetDid,
						MethodSpecificId: "CpeMubv5yw63jXyrgRRsxR",
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
		"cannot get DIDDoc version metadata with an existent old 32 characters Indy style DID, but not existent versionId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.NotExistentIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "notFound",
					DidProperties: types.DidProperties{
						DidString:        testconstants.MigratedIndy32CharStyleTestnetDid,
						MethodSpecificId: "3KpiDD6Hxs4i2G7FtpiGhu",
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
		"cannot get DIDDoc version metadata with an existent DID, but an invalid versionId",
		utils.NegativeTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.TestHostAddress,
				testconstants.UUIDStyleMainnetDid,
				testconstants.InvalidIdentifier,
			),
			ResolutionType: testconstants.DefaultResolutionType,
			ExpectedResult: types.DidResolution{
				Context: "",
				ResolutionMetadata: types.ResolutionMetadata{
					ContentType:     types.JSONLD,
					ResolutionError: "invalidDidUrl",
					DidProperties:   types.DidProperties{},
				},
				Did:      nil,
				Metadata: nil,
			},
			ExpectedStatusCode: errors.InvalidDidUrlHttpCode,
		},
	),
)
