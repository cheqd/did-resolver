package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type DidQueryAllVersionsHandler struct {
	queries.BaseQueryHandler
}

func (d *DidQueryAllVersionsHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	did := service.GetDid()
	contentType := service.GetContentType()

	result, err := c.DidDocService.GetAllDidDocVersionsMetadata(did, contentType)
	if err != nil {
		err.IsDereferencing = d.IsDereferencing
		return nil, err
	}
	content, ok := result.ContentStream.(*types.DereferencedDidVersionsList)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	// Call the next handler
	return d.Continue(c, service, content.Versions)
}
