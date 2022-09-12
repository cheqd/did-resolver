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
	DIDSchemaJSONLD        = "https://www.w3.org/ns/did/v1"
	ResolutionSchemaJSONLD = "https://w3id.org/did-resolution/v1"
)

const (
	DID_METHOD    = "cheqd"
	RESOLVER_PATH = "/1.0/identifiers/"
	RESOURCE_PATH = "/resources/"
)
