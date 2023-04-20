//go:build integration

package query_test

import . "github.com/onsi/ginkgo/v2"

var _ = DescribeTable("Negative: request with common query parameters", func() {

},

	Entry(
		"cannot get DIDDoc with combination of versionId and relativeRef query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceName query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceType query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceVersionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceMetadata query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceCollectionId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and resourceVersion query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of versionId and checksum query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and relativeRef query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceName query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceType query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceVersionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceMetadata query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceCollectionId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and resourceVersion query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of VersionTime and checksum query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and relativeRef query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceName query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceType query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceVersionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceMetadata query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceCollectionId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and resourceVersion query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of metadata and checksum query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and relativeRef query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceName query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceType query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceVersionTime query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceMetadata query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceCollectionId query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and resourceVersion query parameters",
	),

	Entry(
		"cannot get DIDDoc with combination of service and checksum query parameters",
	),
)
