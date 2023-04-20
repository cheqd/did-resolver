//go:build integration

package versionTime

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get DIDDoc with versionTime query", func() {

},

	Entry(
		"can get DIDDoc with versionTime query parameter",
	),

	Entry(
		"can get DIDDoc with an old 16 characters INDY style DID and versionTime query parameter",
	),

	Entry(
		"can get DIDDoc with an old 32 characters INDY style DID and versionTime query parameter",
	),
)
