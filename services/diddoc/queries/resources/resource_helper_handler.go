package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceHelperHandler struct {
}

func (d *ResourceHelperHandler) CastToContent(service services.RequestServiceI, response types.ResolutionResultI) (*types.ResolutionDidDocMetadata, error) {

	// Cast to DidDocMetadataList for getting the list of metadatas
	rc, ok := response.(*types.ResolutionDidDocMetadata)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	return rc, nil
}

// func (d *ResourceHelperHandler) CastToResult(resourceCollectionFiltered types.DereferencedResourceList) *types.ResourceDereferencing {
// 	d.rd.ContentStream = &types.ResolutionDidDocMetadata{
// 		Created:     d.rc.Created,
// 		Updated:     d.rc.Updated,
// 		Deactivated: d.rc.Deactivated,
// 		VersionId:   d.rc.VersionId,
// 		Resources:   resourceCollectionFiltered,
// 	}
// 	return d.rd
// }
