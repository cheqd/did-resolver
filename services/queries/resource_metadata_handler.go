package queries

import (

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceMetadataHandler struct {
	BaseQueryHandler
}

func (d *ResourceMetadataHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	resourceMetadata := service.GetQueryParam(types.ResourceMetadata)
	if resourceMetadata == "" {
		return d.Continue(c, service, response)
	}
	
	// Call the next handler
	return d.Continue(c, service, response)
}
