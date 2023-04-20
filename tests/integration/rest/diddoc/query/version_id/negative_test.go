//go:build integration

package versionId

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Negative: Get DIDDoc with versionId query", func() {

},

	Entry(
		"cannot get DIDDoc with not existent versionId query parameter",
	),

	Entry(
		"cannot get DIDDoc with an invalid versionId query parameter",
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent versionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent service query parameters",
	),

	Entry(
		"cannot get DIDDoc with an existent versionId, but not existent versionTime and service query parameters",
	),
)
