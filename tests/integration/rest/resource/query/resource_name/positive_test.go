//go:build integration

package resource_name_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource with resourceName query", func() {

},

	Entry(
		"can get resource with an existent resourceName query parameter",
	),

	Entry(
		"can get resource with an old 16 characters INDY style DID and resourceName query parameter",
	),

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceName query parameter",
	),
)
