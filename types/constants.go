package types

type ContentType string

const (
	DIDJSON   ContentType = "application/did+json"
	DIDJSONLD ContentType = "application/did+ld+json"
	JSONLD    ContentType = "application/ld+json"
	JSON      ContentType = "application/json"
)

func (cType ContentType) IsSupported() bool {
	supportedTypes := map[ContentType]bool{
		DIDJSON:   true,
		DIDJSONLD: true,
		JSONLD:    true,
	}
	return supportedTypes[cType]
}

const (
	DIDSchemaJSONLD                  = "https://www.w3.org/ns/did/v1"
	ResolutionSchemaJSONLD           = "https://w3id.org/did-resolution/v1"
	Ed25519VerificationKey2020JSONLD = "https://w3id.org/security/suites/ed25519-2020/v1"
	Ed25519VerificationKey2018JSONLD = "https://w3id.org/security/suites/ed25519-2018/v1"
	JsonWebKey2020JSONLD             = "https://w3id.org/security/suites/jws-2020/v1"
)

const (
	DID_METHOD        = "cheqd"
	RESOLVER_PATH     = "/1.0/identifiers/"
	DID_VERSION_PATH  = "/version/"
	DID_VERSIONS_PATH = "/versions"
	DID_METADATA      = "/metadata"
	RESOURCE_PATH     = "/resources/"
	SWAGGER_PATH      = "/swagger/*"
)

const (
	VersionId            string = "versionId"
	VersionTime          string = "versionTime"
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
)
