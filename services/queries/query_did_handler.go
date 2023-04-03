package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DidQueryHandler struct {
	BaseQueryHandler
}

func (d *DidQueryHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	versionId := service.GetQueryParam(types.VersionId)
	versionTime := service.GetQueryParam(types.VersionTime)

	// Here we are handling only query DID without versionId and versionTime
	if versionId != "" || versionTime != "" {
		return d.Continue(c, service, response)
	}
	// Get Params
	did := service.GetDid()
	contentType := service.GetContentType()

	result, err := c.DidDocService.Resolve(did, "", contentType)
	if err != nil {
		err.IsDereferencing = d.IsDereferencing
		return nil, err
	}
	// Call the next handler
	return d.Continue(c, service, result)
}
