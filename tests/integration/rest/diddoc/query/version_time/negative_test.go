//go:build integration

package versionTime

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Negative: Get DIDDoc with versionTime query", func() {

},

	Entry(
		"cannot get DIDDoc with not existent versionTime query parameter",
	),

	Entry(
		"cannot get DIDDoc with an invalid versionTime query parameter",
	),

	Entry(
		"cannot get DIDDoc with an existent versionTime, but not existent versionId query parameters",
	),

	Entry(
		"cannot get DIDDoc with an existent versionTime, but not existent service query parameters",
	),

	Entry(
		"cannot get DIDDoc with an existent versionTime, but not existent versionId and service query parameters",
	),
)
