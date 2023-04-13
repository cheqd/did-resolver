package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type DidDocResolveHandler struct {
	queries.BaseQueryHandler
	DidDocHelperHandler
}

func (dd *DidDocResolveHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	allVersions, err := dd.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	if len(allVersions.Versions) == 0 {
		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, dd.IsDereferencing)
	}

	// Get the latest version. If versionId and versionTime handlers were called here should only 1 element.
	// If versionId or versionTime was not called, we will return the latest version
	versionId := allVersions.Versions[len(allVersions.Versions)-1].VersionId

	result, _err := c.DidDocService.Resolve(service.GetDid(), versionId, service.GetContentType())
	if _err != nil {
		_err.IsDereferencing = dd.IsDereferencing
		return nil, _err
	}

	// Call the next handler
	return dd.Continue(c, service, result)
}
