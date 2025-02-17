package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type ServiceHandler struct {
	queries.BaseQueryHandler
}

func (s *ServiceHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	serviceValue := service.GetQueryParam(types.ServiceQ)

	// If serviceValue is empty, call the next handler. We don't need to handle it here
	if serviceValue == "" {
		return s.Continue(c, service, response)
	}
	// We expect here only DidResolution
	didResolution, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), types.DIDJSONLD, nil, service.GetDereferencing())
	}

	result, err := didResolution.GetServiceByName(serviceValue)
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, types.NewNotFoundError(service.GetDid(), types.DIDJSONLD, nil, service.GetDereferencing())
	}

	// Call the next handler
	return s.Continue(c, service, types.NewServiceResult(result))
}
