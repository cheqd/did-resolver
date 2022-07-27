package types

type ErrorType string

const (
	ResolutionInvalidDID         ErrorType = "invalidDid"
	ResolutionNotFound           ErrorType = "notFound"
	ResolutionMethodNotSupported ErrorType = "methodNotSupported"
)

const (
	DereferencingInvalidDIDUrl ErrorType = "invalidDidUrl"
	DereferencingNotFound      ErrorType = "notFound"
	DereferencingNotSupported  ErrorType = "urlNotSupported"
)

type ContentType string

const (
	DIDJSON   ContentType = "application/did+json"
	DIDJSONLD ContentType = "application/did+ld+json"
	JSONLD    ContentType = "application/ld+json"
)

const (
	DIDSchemaJSONLD = "https://www.w3.org/ns/did/v1"
)

const (
	RESOURCE_PATH = "/resources/"
)
