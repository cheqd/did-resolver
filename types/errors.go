package types

import (
	"errors"
	"fmt"
)

var (
	InvalidDidHttpCode                 = 400
	InvalidDidUrlHttpCode              = 400
	NotFoundHttpCode                   = 404
	RepresentationNotSupportedHttpCode = 406
	InternalErrorHttpCode              = 500
	MethodNotSupportedHttpCode         = 501
)

type IdentityError struct {
	Code            int
	Message         string
	Internal        error
	Did             string
	ContentType     ContentType
	IsDereferencing bool
}

// Error makes it compatible with `error` interface.
func (he *IdentityError) Error() string {
	return fmt.Sprintf("code=%d, message=%s", he.Code, he.Message)
}

func (e IdentityError) GetResolutionOutput() ResolutionResultI {
	metadata := NewResolutionMetadata(e.Did, e.ContentType, e.Message)
	return DidResolution{ResolutionMetadata: metadata}
}

func (e IdentityError) GetDereferencingOutput() ResolutionResultI {
	metadata := NewDereferencingMetadata(e.Did, e.ContentType, e.Message)
	return DidDereferencing{DereferencingMetadata: metadata}
}

func (e *IdentityError) DisplayMessage() ResolutionResultI {
	if e.IsDereferencing {
		return e.GetDereferencingOutput()
	}
	return e.GetResolutionOutput()
}

func NewIdentityError(code int, message string, isDereferencing bool, did string, contentType ContentType, err error) *IdentityError {
	e := IdentityError{
		Code:            code,
		Message:         message,
		Internal:        err,
		Did:             did,
		ContentType:     contentType,
		IsDereferencing: isDereferencing,
	}
	return &e
}

func NewInvalidDIDError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(InvalidDidHttpCode, "invalidDid", isDereferencing, did, contentType, err)
}

func NewInvalidDIDUrlError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(InvalidDidUrlHttpCode, "invalidDidUrl", isDereferencing, did, contentType, err)
}

func NewNotFoundError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(NotFoundHttpCode, "notFound", isDereferencing, did, contentType, err)
}

func NewRepresentationNotSupportedError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(RepresentationNotSupportedHttpCode, "representationNotSupported", isDereferencing, did, contentType, err)
}

func NewInternalError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(InternalErrorHttpCode, "internalError", isDereferencing, did, contentType, err)
}

func NewMethodNotSupportedError(did string, contentType ContentType, err error, isDereferencing bool) *IdentityError {
	return NewIdentityError(MethodNotSupportedHttpCode, "methodNotSupported", isDereferencing, did, contentType, err)
}

func NewInvalidIdentifierError() error {
	return errors.New("unique id should be one of: 16 bytes of decoded base58 string or UUID")
}
