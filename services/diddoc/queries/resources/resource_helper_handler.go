package resources

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type ResourceHelperHandler struct{}

func (d *ResourceHelperHandler) CastToContent(service services.RequestServiceI, response types.ResolutionResultI) (*types.DidResolution, error) {
	// Cast to DidDocMetadataList for getting the list of metadata
	rc, ok := response.(*types.DidResolution)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	return rc, nil
}
