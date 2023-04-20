package diddoc

import (
	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

type DidDocHelperHandler struct{}

func (d *DidDocHelperHandler) CastToContent(service services.RequestServiceI, response types.ResolutionResultI) (types.DidDocMetadataList, error) {
	// Cast to DidDocMetadataList for getting the list of metadata
	rc, ok := response.(types.DidDocMetadataList)
	if !ok {
		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), nil, service.GetDereferencing())
	}
	return rc, nil
}
