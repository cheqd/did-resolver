package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ServiceHandler struct {
	BaseQueryHandler
}

func (s *ServiceHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	serviceValue := service.GetQueryParam(types.ServiceQ)

	// If serviceValue is empty, call the next handler. We don't need to handle it here
	if serviceValue == "" {
		return s.next.Handle(c, service, response)
	}

	// We expect here only DidResolution
	didResolution, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError("response is not DidResolution", types.DIDJSONLD, nil, s.IsDereferencing)
	}

	result, err := didResolution.GetServiceByName(serviceValue)
	if err != nil {
		return nil, err
	}

	// Call the next handler
	return s.Continue(c, service, types.NewServiceResult(result))
}
