package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceVersionHandler struct {
	BaseQueryHandler
}

func (d *ResourceVersionHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// resourceVersion := service.GetQueryParam(types.ResourceVersion)
	
	// Call the next handler
	return d.Continue(c, service, response)
}
