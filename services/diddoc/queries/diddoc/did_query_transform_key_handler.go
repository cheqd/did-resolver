package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type TransformKeyHandler struct {
	queries.BaseQueryHandler
}

func (t *TransformKeyHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	transformKey := types.TransformKeyType(service.GetQueryParam(types.TransformKey))

	// If transformKey is empty, call the next handler. We don't need to handle it here
	if transformKey == "" {
		return t.Continue(c, service, response)
	}

	// We expect here only DidResolution
	didResolution, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError("response is not DidResolution", types.DIDJSONLD, nil, t.IsDereferencing)
	}

	for i, vMethod := range didResolution.Did.VerificationMethod {
		result, err := transformVerificationMethodKey(vMethod, transformKey)
		if err != nil {
			return nil, err
		}
		didResolution.Did.VerificationMethod[i] = result
	}

	// Call the next handler
	return t.Continue(c, service, didResolution)
}
