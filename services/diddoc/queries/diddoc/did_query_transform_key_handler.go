package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type TransformKeysHandler struct {
	queries.BaseQueryHandler
}

func (t *TransformKeysHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	transformKeys := types.TransformKeysType(service.GetQueryParam(types.TransformKeys))

	// If transformKeys is empty, call the next handler. We don't need to handle it here
	if transformKeys == "" {
		return t.Continue(c, service, response)
	}

	if !transformKeys.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(service.GetDid(), types.DIDJSONLD, nil, t.IsDereferencing)
	}

	// We expect here only DidResolution
	didResolution, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), types.DIDJSONLD, nil, t.IsDereferencing)
	}

	for i, vMethod := range didResolution.Did.VerificationMethod {
		result, err := transformVerificationMethodKey(vMethod, transformKeys)
		if err != nil {
			return nil, err
		}
		didResolution.Did.VerificationMethod[i] = result
	}

	// Call the next handler
	return t.Continue(c, service, didResolution)
}
