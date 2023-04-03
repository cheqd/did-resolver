//go:build integration

package rest

import (
	"fmt"
	"net/http"

	testconstants "github.com/cheqd/did-resolver/tests/constants"
	"github.com/go-resty/resty/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = DescribeTable("Test status of redirect DID", func(testCase positiveTestCase) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, err := client.R().
		SetHeader("Accept", testCase.resolutionType).
		Get(testCase.didURL)
	Expect(err).NotTo(BeNil())
	Expect(testCase.expectedStatusCode).To(Equal(resp.StatusCode()))
},

	Entry(
		"can redirect when it try to get DIDDoc with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version metadata with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version metadata with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/version/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc versions with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc versions with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/versions",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource metadata with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource metadata with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get collection of resources with an old 16 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get collection of resources with an old 32 characters Indy style DID",
		positiveTestCase{
			didURL: fmt.Sprintf(
				"http://localhost:8080/1.0/identifiers/%s/metadata",
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			resolutionType:     testconstants.DefaultResolutionType,
			expectedStatusCode: http.StatusMovedPermanently,
		},
	),
)
