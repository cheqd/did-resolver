//go:build integration

package query_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: request with common query parameters", func() {

},

	Entry(
		"can get DIDDoc with existent versionId and versionTime query parameters",
	),

	Entry(
		"can get DIDDoc metadata with existent metadata and versionId query parameters",
	),

	Entry(
		"can get DIDDoc metadata with existent metadata and versionTime query parameters",
	),

	Entry(
		"can get DIDDoc metadata with existent metadata, versionId, versionTime query parameters",
	),

	Entry(
		"can redirect to serviceEndpoint with existent service and versionId query parameters",
	),

	Entry(
		"can redirect to serviceEndpoint with existent service and versionTime query parameters",
	),

	Entry(
		"can redirect to serviceEndpoint with existent service, relativeRef, versionId, versionTime query parameters",
	),
)
