package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DidDocHelperHandler struct {
	rd *types.DidDereferencing
	rc *types.DereferencedDidVersionsList
}

func (d *DidDocHelperHandler) CastToContent(service services.RequestServiceI, response types.ResolutionResultI) (*types.DereferencedDidVersionsList, error) {
	rd, ok := response.(*types.DidDereferencing)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	d.rd = rd

	// Cast to DereferencedResourceListStruct for getting the list of metadatas
	rc, ok := d.rd.ContentStream.(*types.DereferencedDidVersionsList)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	d.rc = rc
	return rc, nil
}

func (d *DidDocHelperHandler) CastToResult(versionFiltered types.DidDocMetadataList) *types.DidDereferencing {
	d.rd.ContentStream = &types.DereferencedDidVersionsList{
		Versions: versionFiltered,
	}
	return d.rd
}
