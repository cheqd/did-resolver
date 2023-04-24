//go:build integration

package resource_checksum_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource with checksum query", func() {

},

	Entry(
		"can get resource with an existent checksum query parameter",
	),

	Entry(
		"can get resource with an old 16 characters INDY style DID and an existent checksum query parameter",
	),

	Entry(
		"can get resource with an old 32 characters INDY style DID and an existent checksum query parameter",
	),
)
