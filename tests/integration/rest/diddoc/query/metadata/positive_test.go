//go:build integration

package metadata

import (
	"encoding/json"
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/integration/rest"

	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get DIDDoc with metadata query", func(testCase utils.PositiveTestCase) {
	client := resty.New()

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).To(BeNil())

	var receivedDidDereferencing utils.DereferencingResult
	Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	var expectedDidDereferencing utils.DereferencingResult
	Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get DIDDoc metadata with metadata=true query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true",
				testconstants.SUTHost,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/metadata/diddoc_metadata_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with metadata=false query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=false",
				testconstants.SUTHost,
				testconstants.IndyStyleMainnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/metadata/diddoc_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an old 16 characters INDY style DID and metadata query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/metadata/diddoc_16_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),

	Entry(
		"can get DIDDoc metadata with an old 32 characters INDY style DID and metadata query parameter",
		utils.PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?metadata=true",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedJSONPath:   "../../../testdata/query/metadata/diddoc_32_old_indy_did.json",
			ExpectedStatusCode: http.StatusOK,
		},
	),
)
