//go:build integration

package resource_metadata_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get Resource Metadata with resourceMetadata query", func() {

},

	Entry(
		"can get resource metadata with metadata=true query parameter",
	),

	Entry(
		"can get resource metadata with metadata=false query parameter",
	),

	Entry(
		"can get collection of resources with an old 16 characters INDY style DID and metadata query parameter",
	),
)
