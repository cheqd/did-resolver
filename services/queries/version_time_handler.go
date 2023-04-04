package queries

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type VersionTimeHandler struct {
	BaseQueryHandler
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

	allMetadatas, err := c.DidDocService.GetAllDidDocVersionsMetadata(did, contentType)
	if err != nil {
		err.IsDereferencing = false
		return nil, err
	}

	allVersions := allMetadatas.ContentStream.(*types.DereferencedDidVersionsList)
	if len(allVersions.Versions) == 0 {
		return nil, types.NewNotFoundError("No versions found", contentType, nil, v.IsDereferencing)
	}

	versionId, _err := allVersions.FindBeforeTime(versionTime)
	if _err != nil {
		return nil, types.NewInternalError("error while finding version before time", contentType, _err, v.IsDereferencing)
	}

	result, err := c.DidDocService.Resolve(did, versionId, contentType)
	if err != nil {
		err.IsDereferencing = v.IsDereferencing
		return nil, err
	}

	// Call the next handler
	return v.Continue(c, service, result)
}
