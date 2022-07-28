package types

type ErrorType string

const (
	InvalidDIDError                 ErrorType = "invalidDid"
	InvalidDIDUrlError              ErrorType = "invalidDidUrl"
	NotFoundError                   ErrorType = "notFound"
	RepresentationNotSupportedError ErrorType = "representationNotSupported"
	MethodNotSupportedError         ErrorType = "methodNotSupported"
	InternalError                   ErrorType = "internalError"
)

func (e ErrorType) GetStatusCode() int {
	switch e {
	case InvalidDIDError:
		return 400
	case InvalidDIDUrlError:
		return 400
	case NotFoundError:
		return 404
	case RepresentationNotSupportedError:
		return 406
	case MethodNotSupportedError:
		return 406
	case "":
		return 200
	default:
		return 500
	}
}

type ContentType string

const (
	DIDJSON   ContentType = "application/did+json"
	DIDJSONLD ContentType = "application/did+ld+json"
	JSONLD    ContentType = "application/ld+json"
	JSON      ContentType = "application/json"
)

func (cType ContentType) IsSupported() bool {
	supportedTypes := map[ContentType]bool{
		DIDJSON: true,
		DIDJSONLD: true,
		JSONLD: true,
	}
	return supportedTypes[cType]
}

const (
	DIDSchemaJSONLD = "https://www.w3.org/ns/did/v1"
)

const (
	RESOURCE_PATH = "/resources/"
)
