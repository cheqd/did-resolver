//go:build integration

package resource_version_time_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Collection of Resources with resourceVersionTime query", func() {

},

	Entry(
		"can get collection of resources with an existent resourceVersionTime query parameter",
	),

	Entry(
		"can get resource with an old 16 characters INDY style DID and resourceVersionTime query parameter",
	),

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceVersionTime query parameter",
	),
)
