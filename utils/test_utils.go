package utils

import (
	"crypto/sha256"

	cheqd "github.com/cheqd/cheqd-node/x/cheqd/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
)

const (
	ValidIdentifier = "N22KY2Dyvmuu2Pyy"
	ValidMethod     = "cheqd"
	ValidNamespace  = "mainnet"
	ValidDid        = "did:" + ValidMethod + ":" + ValidNamespace + ":" + ValidIdentifier
	ValidResourceId = "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
	ValidPubKeyJWK  = "{" +
		"\"crv\":\"Ed25519\"," +
		"\"kid\":\"_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A\"," +
		"\"kty\":\"OKP\"," +
		"\"x\":\"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"" +
		"}"
)

func ValidVerificationMethod() cheqd.VerificationMethod {
	return cheqd.VerificationMethod{
		Id:           ValidDid + "#key-1",
		Type:         "JsonWebKey2020",
		Controller:   ValidDid,
		PublicKeyJwk: cheqd.JSONToPubKeyJWK(ValidPubKeyJWK),
	}
}

func ValidService() cheqd.Service {
	return cheqd.Service{
		Id:              ValidDid + "#service-1",
		Type:            "DIDCommMessaging",
		ServiceEndpoint: "endpoint",
	}
}

func ValidDIDDoc() cheqd.Did {
	service := ValidService()
	verificationMethod := ValidVerificationMethod()

	return cheqd.Did{
		Id:                 ValidDid,
		VerificationMethod: []*cheqd.VerificationMethod{&verificationMethod},
		Service:            []*cheqd.Service{&service},
	}
}

func ValidResource() resource.Resource {
	data := []byte("{\"attr\":[\"name\",\"age\"]}")
	return resource.Resource{
		Header: &resource.ResourceHeader{
			CollectionId: ValidIdentifier,
			Id:           ValidResourceId,
			Name:         "Existing_Resource_Name",
			ResourceType: "CL-Schema",
			MediaType:    "application/json",
			Checksum:     sha256.New().Sum(data),
		},
		Data: data,
	}
}

func ValidMetadata() cheqd.Metadata {
	return cheqd.Metadata{VersionId: "test_version_id", Deactivated: false, Resources: []string{ValidResourceId}}
}

type MockLedgerService struct {
	Did      cheqd.Did
	Metadata cheqd.Metadata
	Resource resource.Resource
}

func NewMockLedgerService(did cheqd.Did, metadata cheqd.Metadata, resource resource.Resource) MockLedgerService {
	return MockLedgerService{
		Did:      did,
		Metadata: metadata,
		Resource: resource,
	}
}

func (ls MockLedgerService) QueryDIDDoc(did string) (*cheqd.Did, *cheqd.Metadata, *types.IdentityError) {
	if did == ls.Did.Id {
		println("query !!!" + ls.Did.Id)
		return &ls.Did, &ls.Metadata, nil
	}
	return nil, nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) QueryResource(did string, resourceId string) (*resource.Resource, *types.IdentityError) {
	if ls.Resource.Header == nil || ls.Resource.Header.Id != resourceId {
		return nil, types.NewNotFoundError(did, types.JSON, nil, true)
	}
	return &ls.Resource, nil
}

func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resource.ResourceHeader, *types.IdentityError) {
	if ls.Metadata.Resources == nil {
		return []*resource.ResourceHeader{}, types.NewNotFoundError(did, types.JSON, nil, true)
	}
	return []*resource.ResourceHeader{ls.Resource.Header}, nil
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}
