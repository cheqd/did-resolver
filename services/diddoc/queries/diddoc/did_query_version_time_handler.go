package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	"github.com/cheqd/did-resolver/types"
)

type VersionTimeHandler struct {
	queries.BaseQueryHandler
	DidDocHelperHandler
}

func (v *VersionTimeHandler) Handle(c services.ResolverContext, service services.RequestServiceI, response types.ResolutionResultI) (types.ResolutionResultI, error) {
	versionTime := service.GetQueryParam(types.VersionTime)

	// Here we are handling only query DID without versionId and versionTime
	if versionTime == "" {
		return v.Continue(c, service, response)
	}
	// Get Params
	did := service.GetDid()
	contentType := service.GetContentType()

	allVersions, err := v.CastToContent(service, response)
	if err != nil {
		return nil, err
	}

	versionId, _err := allVersions.FindActiveForTime(versionTime)
	if _err != nil {
		return nil, types.NewInternalError(did, contentType, _err, v.IsDereferencing)
	}

	if versionId == "" {
		return nil, types.NewNotFoundError(did, contentType, nil, v.IsDereferencing)
	}

	versionsFiltered := allVersions.GetByVersionId(versionId)
	if len(versionsFiltered) == 0 {
		return nil, types.NewInternalError(did, contentType, nil, v.IsDereferencing)
	}

	// Call the next handler
	return v.Continue(c, service, versionsFiltered)
}
