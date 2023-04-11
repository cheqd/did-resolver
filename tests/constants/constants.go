package testconstants

import "fmt"

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
	IndyStyleMainnetDidIdentifier = "4fa8e367-c70e-533e-babf-3732d9761061"
	IndyStyleTestnetDidIdentifier = "60bb3b62-e0f0-545b-a552-63aab5cd1aef"
	UUIDStyleMainnetDidIdentifier = "76e546ee-78cd-5372-b34e-8b47461626e1"
	UUIDStyleTestnetDidIdentifier = "e5615fc2-6f13-42b1-989c-49576a574cef"
)

var (
	UUIDStyleTestnetDidResourceId                    = "9ba3922e-d5f5-4f53-b265-fc0d4e988c77"
	OldIndy32CharStyleTestnetDidIdentifierResourceId = "214b8b61-a861-416b-a7e4-45533af40ada"
)

var (
	NotExistentMainnetDid = fmt.Sprintf(DIDStructure, ValidMethod, ValidMainnetNamespace, NotExistentIdentifier)
	NotExistentTestnetDid = fmt.Sprintf(DIDStructure, ValidMethod, ValidTestnetNamespace, NotExistentIdentifier)
)

var (
	MainnetDIDWithInvalidMethod = fmt.Sprintf(DIDStructure, InvalidMethod, ValidMainnetNamespace, ValidIdentifier)
	TestnetDIDWithInvalidMethod = fmt.Sprintf(DIDStructure, InvalidMethod, ValidTestnetNamespace, ValidIdentifier)

	DIDWithInvalidNamespace = fmt.Sprintf(DIDStructure, ValidMethod, InvalidNamespace, ValidIdentifier)
	InvalidDID              = fmt.Sprintf(DIDStructure, InvalidMethod, InvalidNamespace, InvalidIdentifier)
)

var (
	ValidMethod           = "cheqd"
	ValidMainnetNamespace = "mainnet"
	ValidTestnetNamespace = "testnet"
	ValidIdentifier       = "fb53dd05-329b-4614-a3f2-c0a8c7554ee3"
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

var DIDStructure = "did:%s:%s:%s"

var HashTag = "\u0023"
