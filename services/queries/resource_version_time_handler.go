package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceVersionTimeHandler struct {
	BaseQueryHandler
}

func (d *ResourceVersionTimeHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// resourceVersionTime := service.GetQueryParam(types.ResourceVersionTime)
	
	// Call the next handler
	return d.Continue(c, service, response)
}
