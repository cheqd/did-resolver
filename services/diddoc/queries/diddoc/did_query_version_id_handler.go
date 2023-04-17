package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type VersionIdHandler struct {
	queries.BaseQueryHandler
	DidDocHelperHandler
}

func (v *VersionIdHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	versionId := service.GetQueryParam(types.VersionId)
	// If versionId is empty, call the next handler. We don't need to handle it here
	if versionId == "" {
		return v.Continue(c, service, response)
	}

	// Get Params
	did := service.GetDid()
	contentType := service.GetContentType()
	allVersions, err := v.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	versionFiltered := allVersions.GetByVersionId(versionId)
	if versionFiltered == nil || len(versionFiltered) == 0 {
		return nil, types.NewNotFoundError(did, contentType, nil, service.GetDereferencing())
	}

	versionFiltered[0].Resources = allVersions.GetResourcesBeforeNextVersion(versionId)

	// Call the next handler
	return v.Continue(c, service, versionFiltered)
}
