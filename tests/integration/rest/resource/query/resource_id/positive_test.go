//go:build integration

package resource_id_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource with resourceId query", func() {

},

	Entry(
		"can get resource with an existent resourceId query parameter",
	),

	Entry(
		"can get collection of resources with an old 16 characters INDY style DID and an existent resourceId query parameter",
	),

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and an existent resourceId query parameter",
	),
)
