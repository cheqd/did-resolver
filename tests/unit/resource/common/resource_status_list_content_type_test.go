//go:build unit

package common

import (
	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheqd/did-resolver/services"
	testconstants "github.com/cheqd/did-resolver/tests/constants"
	utils "github.com/cheqd/did-resolver/tests/unit"
	"github.com/cheqd/did-resolver/types"
)

// The ledger sniffs a media type from the content, so a status list token is stored as text/plain or
// application/octet-stream. The resolver serves it with the media type status list clients require.
var _ = Describe("Token Status List content type", func() {
	const (
		jwtId     = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0001"
		cwtId     = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0002"
		taggedId  = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0003"
		otherType = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0004"
		unknownId = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0005"
		jsonId    = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0006"
	)

	jwt := []byte("eyJhbGciOiJFUzI1NiIsInR5cCI6InN0YXR1c2xpc3Qrand0In0.eyJzdWIiOiJ4In0.c2ln")
	// 0x84 (array of four), 0x43 (byte string of three) ...
	cwt := []byte{0x84, 0x43, 0xa1, 0x01, 0x26, 0xa0, 0x41, 0x00, 0x41, 0x00}
	taggedCwt := append([]byte{0xd2}, cwt...)

	resource := func(id string, resourceType string, mediaType string, data []byte) resourceTypes.ResourceWithMetadata {
		return resourceTypes.ResourceWithMetadata{
			Resource: &resourceTypes.Resource{Data: data},
			Metadata: &resourceTypes.Metadata{
				CollectionId: testconstants.ValidIdentifier,
				Id:           id,
				Name:         "revocation-list",
				ResourceType: resourceType,
				MediaType:    mediaType,
			},
		}
	}

	ledger := utils.NewMockLedgerService(
		&testconstants.ValidDIDDoc,
		[]*didTypes.Metadata{&testconstants.ValidMetadata},
		[]resourceTypes.ResourceWithMetadata{
			resource(jwtId, types.ResourceTypeTokenStatusList, "text/plain; charset=utf-8", jwt),
			resource(cwtId, types.ResourceTypeTokenStatusList, "application/octet-stream", cwt),
			resource(taggedId, types.ResourceTypeTokenStatusList, "application/octet-stream", taggedCwt),
			resource(otherType, "SomethingElse", "text/plain; charset=utf-8", jwt),
			resource(unknownId, types.ResourceTypeTokenStatusList, "application/octet-stream", []byte{0x01, 0x02, 0x03}),
			resource(jsonId, types.ResourceTypeTokenStatusList, "application/json", []byte(`{"not":"a token"}`)),
		},
	)
	resourceService := services.NewResourceService(testconstants.ValidMethod, ledger)

	DescribeTable("DereferenceResourceData serves the content type",
		func(resourceId string, expected types.ContentType, expectedBytes []byte) {
			result, err := resourceService.DereferenceResourceData(testconstants.ExistentDid, resourceId, types.StatusListJWT)
			Expect(err).To(BeNil())
			Expect(result.DereferencingMetadata.ContentType).To(Equal(expected))
			Expect(result.ContentStream.GetBytes()).To(Equal(expectedBytes))
		},
		Entry("application/statuslist+jwt for a TokenStatusList resource holding a JWT", jwtId, types.StatusListJWT, jwt),
		Entry("application/statuslist+cwt for a TokenStatusList resource holding a CWT", cwtId, types.StatusListCWT, cwt),
		Entry("application/statuslist+cwt for a tagged CWT", taggedId, types.StatusListCWT, taggedCwt),
		Entry("the ledger media type for a JWT in a resource of another type", otherType, types.ContentType("text/plain; charset=utf-8"), jwt),
		Entry("the ledger media type for a TokenStatusList resource that is neither a JWT nor a CWT", unknownId, types.ContentType("application/octet-stream"), []byte{0x01, 0x02, 0x03}),
		Entry("the ledger media type for a TokenStatusList resource holding JSON", jsonId, types.ContentType("application/json"), []byte(`{"not":"a token"}`)),
	)
})
