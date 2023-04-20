//go:build integration

package metadata

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get DIDDoc with metadata query", func() {

},

	Entry(
		"can get DIDDoc metadata with metadata=true query parameter",
	),

	Entry(
		"can get DIDDoc metadata with metadata=false query parameter",
	),

	Entry(
		"can get DIDDoc metadata with an old 16 characters INDY style DID and metadata query parameter",
	),

	Entry(
		"can get DIDDoc metadata with an old 32 characters INDY style DID and metadata query parameter",
	),
)
