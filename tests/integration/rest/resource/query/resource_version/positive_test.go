//go:build integration

package resource_version_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource with resourceVersion query", func() {

},

	Entry(
		"can get resource with an existent resourceVersion query parameter",
	),

	Entry(
		"can get resource with an old 16 characters INDY style DID and resourceVersion query parameter",
	),

	Entry(
		"can get resource with an old 32 characters INDY style DID and resourceVersion query parameter",
	),
)
