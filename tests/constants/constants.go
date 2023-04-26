package testconstants

import (
	"fmt"
	"time"

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

	UUIDStyleTestnetId  = "c1685ca0-1f5b-439c-8eb8-5c0e85ab7cd0"
	UUIDStyleMainnetDid = "did:cheqd:mainnet:c82f2b02-bdab-4dd7-b833-3e143745d612"
	UUIDStyleTestnetDid = "did:cheqd:testnet:" + UUIDStyleTestnetId

	OldIndy16CharStyleTestnetDid      = "did:cheqd:testnet:zHqbcXb3irKRCMst"
	MigratedIndy16CharStyleTestnetDid = "did:cheqd:testnet:CpeMubv5yw63jXyrgRRsxR"

	OldIndy32CharStyleTestnetDid      = "did:cheqd:testnet:zEv9FXHwp8eFeHbeTXamwda8YoPfgU12"
	MigratedIndy32CharStyleTestnetDid = "did:cheqd:testnet:3KpiDD6Hxs4i2G7FtpiGhu"
)

var (
	IndyStyleMainnetVersionId = "4fa8e367-c70e-533e-babf-3732d9761061"
	IndyStyleTestnetVersionId = "60bb3b62-e0f0-545b-a552-63aab5cd1aef"
	UUIDStyleMainnetVersionId = "76e546ee-78cd-5372-b34e-8b47461626e1"
	UUIDStyleTestnetVersionId = "e5615fc2-6f13-42b1-989c-49576a574cef"
)

var (
	UUIDStyleTestnetDidResourceId                    = "9ba3922e-d5f5-4f53-b265-fc0d4e988c77"
	OldIndy32CharStyleTestnetDidIdentifierResourceId = "214b8b61-a861-416b-a7e4-45533af40ada"
	ExistentResourceName                             = "Demo Resource"
	ExistentResourceType                             = "String"
	ExistentResourceVersion                          = ""
	ExistentResourceMediaType                        = "application/json"
	ExistentResourceChecksum                         = "e1dbc03b50bdb995961dc8843df6539b79d03bf49787ed6462189ee97d27eaf3"
	ExistentResourceCreated                          = "2023-01-25T12:08:39.63Z"
	ExistentResourceVersionTimeAfter                 = "2023-01-26T12:08:39.63Z"
	ExistentResourceVersionTimeBefore                = "2023-01-24T12:08:39.63Z"
	ExistentResource                                 = types.DereferencedResource{
		ResourceURI:  UUIDStyleTestnetId + types.RESOURCE_PATH + UUIDStyleTestnetDidResourceId,
		CollectionId: UUIDStyleTestnetId,
		ResourceId:   UUIDStyleTestnetDidResourceId,
		Name:         ExistentResourceName,
	}
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
	ValidVersionId        = "32e0613e-bee4-4ea4-952c-bba3e857fa2a"
	ValidServiceId        = "service-1"
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
	NotExistentService    = "not_existent_service"
)

var (
	InvalidMethod     = "invalid_method"
	InvalidNamespace  = "invalid_namespace"
	InvalidIdentifier = "invalid_identifier"
	InvalidVersionId  = "invalid_uuid_identifier"
	InvalidServiceId  = "not_found_service_id"
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
	NotEmptyTime     = NotEmptyTimestamp.AsTime()
	ValidCreated, _  = time.Parse(time.RFC3339, "2021-08-23T09:00:00Z")
	CreatedAfter, _  = time.Parse(time.RFC3339, "2021-08-23T09:10:00Z")
	CreatedBefore, _ = time.Parse(time.RFC3339, "2021-08-23T08:00:00Z")
	ValidUpdated, _  = time.Parse(time.RFC3339, "2021-08-23T09:30:00Z")
	UpdatedAfter, _  = time.Parse(time.RFC3339, "2021-08-23T09:30:01Z")
	UpdatedBefore, _ = time.Parse(time.RFC3339, "2021-08-23T09:20:00Z")
)

var (
	ValidDIDDoc                   = generateDIDDoc()
	ValidMetadata                 = generateMetadata()
	ValidResource                 = generateResource()
	ValidVerificationMethod       = generateVerificationMethod()
	ValidService                  = generateService()
	ValidDIDDocResolution         = types.NewDidDoc(&ValidDIDDoc)
	ValidFragmentMetadata         = types.NewResolutionDidDocMetadata(ExistentDid, &ValidMetadata, []*resourceTypes.Metadata{})
	ValidResourceDereferencing    = types.DereferencedResourceData(ValidResource[0].Resource.Data)
	ValidDereferencedResourceList = types.NewDereferencedResourceListStruct(ExistentDid, []*resourceTypes.Metadata{ValidResource[0].Metadata})
	ValidDid                      = ValidDIDDoc.Id
)

var DIDStructure = "did:%s:%s:%s"

var HashTag = "\u0023"
