package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceHelperHandler struct {
	rd *types.ResourceDereferencing
	rc *types.ResolutionDidDocMetadata
}

func (d *ResourceHelperHandler) CastToContent(service services.RequestServiceI, response types.ResolutionResultI) (*types.ResolutionDidDocMetadata, error) {
	rd, ok := response.(*types.ResourceDereferencing)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	d.rd = rd

	// Cast to DereferencedResourceListStruct for getting the list of metadatas
	rc, ok := d.rd.ContentStream.(*types.ResolutionDidDocMetadata)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	d.rc = rc
	return rc, nil
}

func (d *ResourceHelperHandler) CastToResult(resourceCollectionFiltered types.DereferencedResourceList) *types.ResourceDereferencing {
	d.rd.ContentStream = &types.ResolutionDidDocMetadata{
		Created:     d.rc.Created,
		Updated:     d.rc.Updated,
		Deactivated: d.rc.Deactivated,
		VersionId:   d.rc.VersionId,
		Resources:   resourceCollectionFiltered,
	}
	return d.rd
}
