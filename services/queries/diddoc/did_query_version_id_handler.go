package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	"github.com/cheqd/did-resolver/types"
)

type VersionIdHandler struct {
	queries.BaseQueryHandler
}

func (v *VersionIdHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	// Get Params
	did := service.GetDid()
	contentType := service.GetContentType()
	versionId := service.GetQueryParam(types.VersionId)

	// If versionId is empty, call the next handler. We don't need to handle it here
	if versionId == "" {
		return v.Continue(c, service, response)
	}
	result, err := c.DidDocService.Resolve(did, versionId, contentType)
	if err != nil {
		err.IsDereferencing = v.IsDereferencing
		return nil, err
	}
	// Call the next handler
	return v.Continue(c, service, result)
}
