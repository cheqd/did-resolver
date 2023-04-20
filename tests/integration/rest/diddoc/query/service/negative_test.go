//go:build integration

package service

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Negative: Get DIDDoc with service query", func() {

},

	Entry(
		"cannot redirect to serviceEndpoint with not existent service query parameter",
	),

	Entry(
		"cannot redirect to serviceEndpoint with relativeRef query parameter",
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but an invalid relativeRef URI parameters",
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionId query parameters",
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionTime query parameters",
	),

	Entry(
		"cannot redirect to serviceEndpoint with an existent service, but not existent versionId and versionTime query parameters",
	),
)
