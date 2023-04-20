//go:build integration

package versionId

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get DIDDoc with versionId query", func() {

},

	Entry(
		"can get DIDDoc with versionId query parameter",
	),

	Entry(
		"can get DIDDoc with an old 16 characters INDY style DID and versionId query parameter",
	),

	Entry(
		"can get DIDDoc with an old 32 characters INDY style DID and versionId query parameter",
	),
)
