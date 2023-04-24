//go:build integration

package resource_collection_id_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Collection of Resources with collectionId query", func() {

},

	Entry(
		"can get collection of resources with an existent collectionId query parameter",
	),

	Entry(
		"can get collection of resources with an old 16 characters INDY style DID and an existent collectionId query parameter",
	),

	Entry(
		"can get collection of resources with an old 32 characters INDY style DID and an existent collectionId query parameter",
	),
)
