package utils

import (
	"crypto/sha256"
	"fmt"

	didTypes "github.com/cheqd/cheqd-node/x/did/types"
	resource "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
)

const (
	ValidIdentifier = "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
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

func ValidVerificationMethod() didTypes.VerificationMethod {
	return didTypes.VerificationMethod{
		Id:                   ValidDid + "#key-1",
		Type:                 "JsonWebKey2020",
		Controller:           ValidDid,
		VerificationMaterial: ValidPubKeyJWK,
	}
}

func ValidService() didTypes.Service {
	return didTypes.Service{
		Id:              ValidDid + "#service-1",
		Type:            "DIDCommMessaging",
		ServiceEndpoint: []string{"http://example.com"},
	}
}

func ValidDIDDoc() didTypes.DidDoc {
	service := ValidService()
	verificationMethod := ValidVerificationMethod()

	return didTypes.DidDoc{
		Id:                 ValidDid,
		VerificationMethod: []*didTypes.VerificationMethod{&verificationMethod},
		Service:            []*didTypes.Service{&service},
	}
}

func ValidResource() resource.ResourceWithMetadata {
	data := []byte("{\"attr\":[\"name\",\"age\"]}")
	checksum := sha256.New().Sum(data)
	return resource.ResourceWithMetadata{
		Resource: &resource.Resource{
			Data: data,
		},
		Metadata: &resource.Metadata{
			CollectionId: ValidIdentifier,
			Id:           ValidResourceId,
			Name:         ValidResourceId,
			ResourceType: "string",
			MediaType:    "application/json",
			Checksum:     fmt.Sprintf("%x", checksum),
		},
	}
}

func ValidMetadata() didTypes.Metadata {
	return didTypes.Metadata{VersionId: "test_version_id", Deactivated: false}
}

type MockLedgerService struct {
	Did      didTypes.DidDoc
	Metadata didTypes.Metadata
	Resource resource.ResourceWithMetadata
}

func NewMockLedgerService(did didTypes.DidDoc, metadata didTypes.Metadata, resource resource.ResourceWithMetadata) MockLedgerService {
	return MockLedgerService{
		Did:      did,
		Metadata: metadata,
		Resource: resource,
	}
}

func (ls MockLedgerService) QueryDIDDoc(did string) (*didTypes.DidDocWithMetadata, *types.IdentityError) {
	if did == ls.Did.Id {
		println("query !!!" + ls.Did.Id)
		return &didTypes.DidDocWithMetadata{DidDoc: &ls.Did, Metadata: &ls.Metadata}, nil
	}
	return nil, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) QueryResource(did string, resourceId string) (*resource.ResourceWithMetadata, *types.IdentityError) {
	if ls.Resource.Metadata == nil || ls.Resource.Metadata.Id != resourceId {
		return nil, types.NewNotFoundError(did, types.JSON, nil, true)
	}
	return &ls.Resource, nil
}

func (ls MockLedgerService) QueryCollectionResources(did string) ([]*resource.Metadata, *types.IdentityError) {
	return []*resource.Metadata{}, types.NewNotFoundError(did, types.JSON, nil, true)
}

func (ls MockLedgerService) GetNamespaces() []string {
	return []string{"testnet", "mainnet"}
}
