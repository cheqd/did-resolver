package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type StopHandler struct {
	BaseQueryHandler
}

func (s *StopHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {	
	return response, nil
}