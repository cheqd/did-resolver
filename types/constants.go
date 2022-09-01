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
	DIDSchemaJSONLD = "https://www.w3.org/ns/did/v1"
)

const (
	RESOURCE_PATH = "/resources/"
)
