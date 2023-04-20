//go:build integration

package metadata

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Negative: Get DIDDoc with metadata query", func() {

},

	Entry(
		"cannot get DIDDoc metadata with not supported metadata query parameter value",
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionId query parameters",
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent service query parameters",
	),

	Entry(
		"cannot get DIDDoc metadata with supported metadata value, but not existent versionId, versionTime, service query parameters",
	),
)
