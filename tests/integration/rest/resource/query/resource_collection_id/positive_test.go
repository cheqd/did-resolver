//go:build integration

package resource_collection_id_test

import (
	. "github.com/onsi/ginkgo/v2"
	// . "github.com/onsi/gomega"
)

var _ = DescribeTable("Positive: Get Collection of Resources with collectionId query", func() {
	// client := resty.New()

	// resp, err := client.R().
	// 	SetHeader("Accept", testCase.ResolutionType).
	// 	SetHeader("Accept-Encoding", testCase.EncodingType).
	// 	Get(testCase.DidURL)
	// Expect(err).To(BeNil())

	// var receivedDidDereferencing utils.DereferencingResult
	// Expect(json.Unmarshal(resp.Body(), &receivedDidDereferencing)).To(BeNil())
	// Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))

	// var expectedDidDereferencing utils.DereferencingResult
	// Expect(utils.ConvertJsonFileToType(testCase.ExpectedJSONPath, &expectedDidDereferencing)).To(BeNil())

	// Expect(testCase.ExpectedEncodingType).To(Equal(resp.Header().Get("Content-Encoding")))
	// utils.AssertDidDereferencing(expectedDidDereferencing, receivedDidDereferencing)
},

	Entry(
		"can get collection of resources with an existent collectionId query parameter",
		// utils.PositiveTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?collectionId=%s",
		// 		testconstants.UUIDStyleTestnetDid,
		// 		testconstants.UUIDStyleTestnetId,
		// 	),
		// 	ResolutionType:     testconstants.DefaultResolutionType,
		// 	ExpectedJSONPath:   "../../../testdata/query/collection_id/metadata_did.json",
		// 	ExpectedStatusCode: http.StatusOK,
		// },
	),

	// TODO: add unit test for testing get resource with an old 16 characters INDY style DID
	// and collectionId query parameter.

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and an existent collectionId query parameter",
		// utils.PositiveTestCase{
		// 	DidURL: fmt.Sprintf(
		// 		"http://localhost:8080/1.0/identifiers/%s?collectionId=%s",
		// 		testconstants.OldIndy32CharStyleTestnetDid,
		// 		"zEv9FXHwp8eFeHbeTXamwda8YoPfgU12",
		// 	),
		// 	ResolutionType:     testconstants.DefaultResolutionType,
		// 	ExpectedJSONPath:   "../../../testdata/query/collection_id/metadata_32_indy_did.json",
		// 	ExpectedStatusCode: http.StatusOK,
		// },
	),
)
