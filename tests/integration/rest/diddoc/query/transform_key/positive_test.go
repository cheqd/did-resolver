//go:build integration

package transformKey

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"
	"github.com/cheqd/did-resolver/types"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	DidWithEd25519VerificationKey2018Key = "did:cheqd:testnet:d8ac0372-0d4b-413e-8ef5-8e8f07822b2c"
	DidWithJsonWebKey2020Key             = "did:cheqd:testnet:54c96733-32ad-4878-b7ce-f62f4fdf3291"
)

var _ = DescribeTable("Positive: Get DIDDoc with transformKey query parameter", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidResolution types.DidResolution
	Expect(json.Unmarshal(resp.Body(), &receivedDidResolution)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidResolution types.DidResolution
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidResolution)).To(BeNil())

	utils.AssertDidResolution(expectedDidResolution, receivedDidResolution)
},
	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2018 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported JSONWebKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithEd25519VerificationKey2018Key,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported Ed25519VerificationKey2018 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported Ed25519VerificationKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported JSONWebKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2018 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithJsonWebKey2020Key,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithJsonWebKey2020Key,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported JSONWebKey2020 transformKey query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKey=%s",
				DidWithJsonWebKey2020Key,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_jwk_2020.json",
		},
	),
)
