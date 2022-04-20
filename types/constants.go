package types

type ErrorType string

const (
	ResolutionInvalidDID         ErrorType = "invalidDid"
	ResolutionNotFound           ErrorType = "notFound"
	ResolutionMethodNotSupported ErrorType = "methodNotSupported"
)

const (
	DereferencingInvalidDIDUrl    ErrorType = "invalidDidUrl"
	DereferencingFragmentNotFound ErrorType = "FragmentNotFound"
	DereferencingNotSupported     ErrorType = "UrlNotSupported"
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
