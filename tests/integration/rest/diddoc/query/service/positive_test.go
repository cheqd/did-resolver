//go:build integration

package service

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Positive: Get DIDDoc with service query", func() {

},

	Entry(
		"can redirect to serviceEndpoint with an existent service query parameter",
	),

	Entry(
		"can redirect to serviceEndpoint with an existent service and a valid relativeRef URI query parameters",
	),

	Entry(
		"can redirect to serviceEndpoint with an old 16 characters INDY style DID and service query parameter",
	),

	Entry(
		"can redirect to serviceEndpoint with an old 32 characters INDY style DID and service query parameter",
	),
)
