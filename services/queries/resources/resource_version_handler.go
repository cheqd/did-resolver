package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/services/queries"
)

type ResourceVersionHandler struct {
	queries.BaseQueryHandler
}

func (d *ResourceVersionHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// resourceVersion := service.GetQueryParam(types.ResourceVersion)
	
	// Call the next handler
	return d.Continue(c, service, response)
}
