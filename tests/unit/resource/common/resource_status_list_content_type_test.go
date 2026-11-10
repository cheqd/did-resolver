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
		jwtId   = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0001"
		cwtId   = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0002"
		plainId = "4e0f3d6e-3f0a-4a49-9c43-0d4c1c6b0003"
	)

	resource := func(id string, resourceType string, mediaType string) resourceTypes.ResourceWithMetadata {
		return resourceTypes.ResourceWithMetadata{
			Resource: &resourceTypes.Resource{Data: []byte("token")},
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
			resource(jwtId, types.ResourceTypeTokenStatusListJWT, "text/plain; charset=utf-8"),
			resource(cwtId, types.ResourceTypeTokenStatusListCWT, "application/octet-stream"),
			resource(plainId, "SomethingElse", "text/plain; charset=utf-8"),
		},
	)
	resourceService := services.NewResourceService(testconstants.ValidMethod, ledger)

	DescribeTable("DereferenceResourceData serves the content type",
		func(resourceId string, expected types.ContentType) {
			result, err := resourceService.DereferenceResourceData(testconstants.ExistentDid, resourceId, types.StatusListJWT)
			Expect(err).To(BeNil())
			Expect(result.DereferencingMetadata.ContentType).To(Equal(expected))
			Expect(result.ContentStream.GetBytes()).To(Equal([]byte("token")))
		},
		Entry("application/statuslist+jwt for a TokenStatusListJwt resource", jwtId, types.StatusListJWT),
		Entry("application/statuslist+cwt for a TokenStatusListCwt resource", cwtId, types.StatusListCWT),
		Entry("the ledger media type for any other resource", plainId, types.ContentType("text/plain; charset=utf-8")),
	)
})
