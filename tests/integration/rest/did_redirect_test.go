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

var _ = DescribeTable("Test HTTP status code of redirect DID URL", func(testCase PositiveTestCase) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())

	resp, err := client.R().
		SetHeader("Accept", testCase.ResolutionType).
		Get(testCase.DidURL)
	Expect(err).NotTo(BeNil())
	Expect(testCase.ExpectedStatusCode).To(Equal(resp.StatusCode()))
},

	Entry(
		"can redirect when it try to get DIDDoc with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version metadata with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc version metadata with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/version/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc versions with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get DIDDoc versions with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/versions",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource data with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource metadata with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get resource metadata with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/resources/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get collection of resources with an old 16 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),

	Entry(
		"can redirect when it try to get collection of resources with an old 32 characters Indy style DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s/metadata",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),
	Entry(
		"can redirect when it try to get resource by query params. 16 symbols DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceId=%s",
				testconstants.SUTHost,
				testconstants.OldIndy16CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),
	Entry(
		"can redirect when it try to get resource by query params. 32 symbols DID",
		PositiveTestCase{
			DidURL: fmt.Sprintf(
				"http://%s/1.0/identifiers/%s?resourceId=%s",
				testconstants.SUTHost,
				testconstants.OldIndy32CharStyleTestnetDid,
				testconstants.ValidIdentifier,
			),
			ResolutionType:     testconstants.DefaultResolutionType,
			ExpectedStatusCode: http.StatusMovedPermanently,
		},
	),
)
