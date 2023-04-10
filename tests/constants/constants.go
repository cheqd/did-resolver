package testconstants

import (
	"fmt"

	resourceTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/resource/v2"
	"github.com/cheqd/did-resolver/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	DefaultResolutionType    = "*/*"
	DefaultEncodingType      = "gzip, deflate, br"
	NotSupportedEncodingType = "deflate, br"
)

var (
	IndyStyleMainnetDid = "did:cheqd:mainnet:Ps1ysXP2Ae6GBfxNhNQNKN"
	IndyStyleTestnetDid = "did:cheqd:testnet:73wnEyHhkhXiH1Nq7w5Kgq"

	UUIDStyleMainnetDid = "did:cheqd:mainnet:c82f2b02-bdab-4dd7-b833-3e143745d612"
	UUIDStyleTestnetDid = "did:cheqd:testnet:c1685ca0-1f5b-439c-8eb8-5c0e85ab7cd0"

	OldIndy16CharStyleTestnetDid      = "did:cheqd:testnet:zHqbcXb3irKRCMst"
	MigratedIndy16CharStyleTestnetDid = "did:cheqd:testnet:CpeMubv5yw63jXyrgRRsxR"

	OldIndy32CharStyleTestnetDid      = "did:cheqd:testnet:zEv9FXHwp8eFeHbeTXamwda8YoPfgU12"
	MigratedIndy32CharStyleTestnetDid = "did:cheqd:testnet:3KpiDD6Hxs4i2G7FtpiGhu"
)

var (
	NotExistentMainnetDid = fmt.Sprintf(DIDStructure, ValidMethod, ValidMainnetNamespace, NotExistentIdentifier)
	NotExistentTestnetDid = fmt.Sprintf(DIDStructure, ValidMethod, ValidTestnetNamespace, NotExistentIdentifier)
)

var (
	MainnetDidWithInvalidMethod = fmt.Sprintf(DIDStructure, InvalidMethod, ValidMainnetNamespace, ValidIdentifier)
	TestnetDidWithInvalidMethod = fmt.Sprintf(DIDStructure, InvalidMethod, ValidTestnetNamespace, ValidIdentifier)

	DidWithInvalidNamespace = fmt.Sprintf(DIDStructure, ValidMethod, InvalidNamespace, ValidIdentifier)
	InvalidDid              = fmt.Sprintf(DIDStructure, InvalidMethod, InvalidNamespace, InvalidIdentifier)
)

var (
	ValidMethod           = "cheqd"
	ValidMainnetNamespace = "mainnet"
	ValidTestnetNamespace = "testnet"
	ValidIdentifier       = "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
	ValidVersionId        = "valid_version_id"
	ValidPubKeyJWK        = "{" +
		"\"crv\":\"Ed25519\"," +
		"\"kid\":\"_Qq0UL2Fq651Q0Fjd6TvnYE-faHiOpRlPVQcY_-tA4A\"," +
		"\"kty\":\"OKP\"," +
		"\"x\":\"VCpo2LMLhn6iWku8MKvSLg2ZAoC-nlOyPVQaO3FxVeQ\"" +
		"}"
)

var (
	ExistentDid        = fmt.Sprintf(DIDStructure, ValidMethod, ValidMainnetNamespace, ValidIdentifier)
	ExistentResourceId = "a09abea0-22e0-4b35-8f70-9cc3a6d0b5fd"
)

var (
	ValidResourceData     = []byte("test_checksum")
	ValidResourceMetadata = resourceTypes.Metadata{
		CollectionId: ValidIdentifier,
		Id:           ExistentResourceId,
		Name:         "Existing Resource Name",
		ResourceType: "CL-Schema",
		MediaType:    "application/json",
		Checksum:     generateChecksum(ValidResourceData),
	}

	ValidMetadataResource = types.DereferencedResource{
		ResourceURI:       ExistentDid + types.RESOURCE_PATH + ValidResourceMetadata.Id,
		CollectionId:      ValidResourceMetadata.CollectionId,
		ResourceId:        ValidResourceMetadata.Id,
		Name:              ValidResourceMetadata.Name,
		ResourceType:      ValidResourceMetadata.ResourceType,
		MediaType:         ValidResourceMetadata.MediaType,
		Created:           &EmptyTime,
		Checksum:          ValidResourceMetadata.Checksum,
		PreviousVersionId: nil,
		NextVersionId:     nil,
	}
)

var (
	NotExistentIdentifier = "ffffffff-329b-4614-a3f2-ffffffffffff"
	NotExistentFragment   = "not_existent_fragment"
)

var (
	InvalidMethod     = "invalid_method"
	InvalidNamespace  = "invalid_namespace"
	InvalidIdentifier = "invalid_identifier"
)

var (
	EmptyTimestamp = &timestamppb.Timestamp{
		Seconds: 0,
		Nanos:   0,
	}
	EmptyTime = EmptyTimestamp.AsTime()

	NotEmptyTimestamp = &timestamppb.Timestamp{
		Seconds: 123456789,
		Nanos:   0,
	}
	NotEmptyTime = NotEmptyTimestamp.AsTime()
)

var (
	ValidDIDDoc                   = generateDIDDoc()
	ValidMetadata                 = generateMetadata()
	ValidResource                 = generateResource()
	ValidVerificationMethod       = generateVerificationMethod()
	ValidService                  = generateService()
	ValidDIDDocResolution         = types.NewDidDoc(&ValidDIDDoc)
	ValidFragmentMetadata         = types.NewResolutionDidDocMetadata(ExistentDid, &ValidMetadata, []*resourceTypes.Metadata{})
	ValidResourceDereferencing    = types.DereferencedResourceData(ValidResource.Resource.Data)
	ValidDereferencedResourceList = types.NewDereferencedResourceList(ExistentDid, []*resourceTypes.Metadata{ValidResource.Metadata})
)

var DIDStructure = "did:%s:%s:%s"

var HashTag = "\u0023"
