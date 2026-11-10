package types

import (
	"bytes"
	"regexp"
)

type ContentType string

const (
	DIDJSON    ContentType = "application/did+json"
	DIDJSONLD  ContentType = "application/did+ld+json"
	JSONLD     ContentType = "application/ld+json"
	JSON       ContentType = "application/json"
	DIDRES     ContentType = "application/did"
	W3IDDIDRES string      = "https://w3id.org/did-resolution"
	TEXT       ContentType = "text/plain"
	// Media types of Token Status List tokens (draft-ietf-oauth-status-list). The ledger sniffs a media type
	// from the content, so these cannot be recorded on chain and are derived from the resource type.
	StatusListJWT ContentType = "application/statuslist+jwt"
	StatusListCWT ContentType = "application/statuslist+cwt"
	W3IDDIDURL    string      = "https://w3id.org/did-url-dereferencing"
)

func (cType ContentType) IsSupported() bool {
	supportedTypes := map[ContentType]bool{
		DIDJSON:   true,
		DIDJSONLD: true,
		JSONLD:    true,
		DIDRES:    true,
		JSON:      true,
	}
	return supportedTypes[cType]
}

type TransformKeysType string

const (
	Ed25519VerificationKey2018 TransformKeysType = "Ed25519VerificationKey2018"
	Ed25519VerificationKey2020 TransformKeysType = "Ed25519VerificationKey2020"
	JsonWebKey2020             TransformKeysType = "JsonWebKey2020"
)

func (tKType TransformKeysType) IsSupported() bool {
	supportedTypes := map[TransformKeysType]bool{
		Ed25519VerificationKey2018: true,
		Ed25519VerificationKey2020: true,
		JsonWebKey2020:             true,
	}
	return supportedTypes[tKType]
}

const (
	DIDSchemaJSONLD                  = "https://www.w3.org/ns/did/v1"
	LinkedDomainsJSONLD              = "https://identity.foundation/.well-known/did-configuration/v1"
	ResolutionSchemaJSONLD           = "https://w3id.org/did-resolution/v1"
	Ed25519VerificationKey2020JSONLD = "https://w3id.org/security/suites/ed25519-2020/v1"
	Ed25519VerificationKey2018JSONLD = "https://w3id.org/security/suites/ed25519-2018/v1"
	JsonWebKey2020JSONLD             = "https://w3id.org/security/suites/jws-2020/v1"
)

const (
	DID_METHOD              = "cheqd"
	RESOLVER_PATH           = "/1.0/identifiers/"
	DID_VERSION_PATH        = "/version/"
	DID_VERSIONS_PATH       = "/versions"
	DID_METADATA            = "/metadata"
	RESOURCE_PATH           = "/resources/"
	SWAGGER_PATH            = "/swagger/*"
	DEFAULT_RESOLUTION_TYPE = "*/*"
)

const (
	VersionId            string = "versionId"
	VersionTime          string = "versionTime"
	TransformKeys        string = "transformKeys"
	LinkedDomains        string = "LinkedDomains"
	Metadata             string = "metadata"
	ServiceQ             string = "service"
	RelativeRef          string = "relativeRef"
	ResourceId           string = "resourceId"
	ResourceName         string = "resourceName"
	ResourceType         string = "resourceType"
	ResourceVersionTime  string = "resourceVersionTime"
	ResourceMetadata     string = "resourceMetadata"
	ResourceCollectionId string = "resourceCollectionId"
	ResourceVersion      string = "resourceVersion"
	ResourceChecksum     string = "checksum"
)

// ResourceTypeTokenStatusList is the resource type cheqd Studio documents for Token Status Lists. A resource of this
// type holds either a JWT or a CWT, so the token format is recognised from the content.
const ResourceTypeTokenStatusList = "TokenStatusList"

var jwtPattern = regexp.MustCompile(`^[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]*$`)

// ResourceContentType returns the content type to serve a resource with: the Token Status List media type for a
// status list token, and the media type recorded on the ledger for everything else. The ledger sniffs a media type from
// the content, so it records a status list as text/plain or application/octet-stream. Content that is not recognised as
// a JWT or CWT keeps the ledger media type, so nothing is mislabelled.
func ResourceContentType(resourceType string, ledgerMediaType string, data []byte) ContentType {
	if resourceType != ResourceTypeTokenStatusList {
		return ContentType(ledgerMediaType)
	}

	if isCoseSign1(data) {
		return StatusListCWT
	}
	if jwtPattern.Match(bytes.TrimSpace(data)) {
		return StatusListJWT
	}
	return ContentType(ledgerMediaType)
}

// isCoseSign1 reports whether data starts like a COSE_Sign1: the four element CBOR array 0x84, optionally preceded by
// tag 18 (0xD2), followed by the protected header as a byte string (major type 2, 0x40 to 0x5F).
func isCoseSign1(data []byte) bool {
	i := 0
	if len(data) > 0 && data[0] == 0xD2 {
		i = 1
	}
	return len(data) > i+1 && data[i] == 0x84 && data[i+1] >= 0x40 && data[i+1] <= 0x5F
}
