package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type RelativeRefHandler struct {
	BaseQueryHandler
}

func (r *RelativeRefHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	relativeRef := service.GetQueryParam(types.RelativeRef)

	// If relativeRef is empty, call the next handler. We don't need to handle it here
	if relativeRef == "" {
		return r.Continue(c, service, response)
	}

	// We expect here only DidResolution
	serviceResult, ok := response.(*types.ServiceResult)
	if !ok {
		return r.Continue(c, service, response)
	}

	// Call the next handler
	return r.Continue(c, service, types.NewServiceResult(serviceResult.GetServiceEndpoint()+relativeRef))
}
