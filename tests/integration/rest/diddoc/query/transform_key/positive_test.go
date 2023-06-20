//go:build integration

package transformKeys

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
	didWithEd25519VerificationKey2018Key = "did:cheqd:testnet:d8ac0372-0d4b-413e-8ef5-8e8f07822b2c"
	didWithJsonWebKey2020Key             = "did:cheqd:testnet:54c96733-32ad-4878-b7ce-f62f4fdf3291"
)

var _ = DescribeTable("Positive: Get DIDDoc with transformKeys query parameter", func(testCase utils.PositiveTestCase) {
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
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2018 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported JSONWebKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithEd25519VerificationKey2018Key,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2018_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported Ed25519VerificationKey2018 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported Ed25519VerificationKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with supported JSONWebKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.UUIDStyleTestnetDid,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_ed25519_2020_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2018 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithJsonWebKey2020Key,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported Ed25519VerificationKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithJsonWebKey2020Key,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with supported JSONWebKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				didWithJsonWebKey2020Key,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_jwk_2020_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with an existent old 16 characters INDY style DID and supported Ed25519VerificationKey2018 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_16_indy_jwk_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with an existent old 16 characters INDY style DID and supported Ed25519VerificationKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_16_indy_jwk_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (JSONWebKey2020) with an existent old 16 characters INDY style DID and supported JSONWebKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_16_indy_jwk_2020_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with an existent old 32 characters INDY style DID and supported Ed25519VerificationKey2018 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				string(types.Ed25519VerificationKey2018),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_32_indy_ed255519_2020_to_ed25519_2018.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with an existent old 32 characters INDY style DID and supported Ed25519VerificationKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				string(types.Ed25519VerificationKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_32_indy_ed255519_2020_to_ed25519_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2020) with an existent old 32 characters INDY style DID and supported JSONWebKey2020 transformKeys query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				string(types.JsonWebKey2020),
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_old_32_indy_ed255519_2020_to_jwk_2020.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2020 transformKeys and DIDDoc versionId queries",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s&versionId=%s",
				didWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2020),
				"44f49254-8106-40ee-99ad-e50ac9517346",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_transform_key_and_version_id.json",
		},
	),

	Entry(
		"can get DIDDoc (Ed25519VerificationKey2018) with supported Ed25519VerificationKey2020 transformKeys and DIDDoc versionTime queries",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s?transformKeys=%s&versionTime=%s",
				didWithEd25519VerificationKey2018Key,
				string(types.Ed25519VerificationKey2020),
				"2023-02-21T14:28:48.406713879Z",
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusOK,
			ExpectedJSONPath:   "../../../testdata/query/transform_key/diddoc_transform_key_version_time.json",
		},
	),
)
