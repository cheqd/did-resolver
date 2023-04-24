//go:build integration

package resource_type_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource with resourceType query", func() {

},

	Entry(
		"can get resource with an existent resourceType query parameter",
	),

	Entry(
		"can get resource with an old 16 characters INDY style DID and resourceType query parameter",
	),

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceType query parameter",
	),
)
