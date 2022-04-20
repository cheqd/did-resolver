package types

type ResolutionError string

const (
	ResolutionInvalidDID         ResolutionError = "invalidDid"
	ResolutionNotFound           ResolutionError = "notFound"
	ResolutionMethodNotSupported ResolutionError = "methodNotSupported"
)

type ContentType string

const (
	DIDJSON   ContentType = "application/did+json"
	DIDJSONLD ContentType = "application/did+ld+json"
	JSONLD    ContentType = "application/ld+json"
)

const (
	DIDSchemaJSONLD = "https://ww.w3.org/ns/did/v1"
)
