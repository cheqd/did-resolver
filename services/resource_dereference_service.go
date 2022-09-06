package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type ResourceService struct {
	didMethod     string
	ledgerService LedgerServiceI
}

func NewResourceService(didMethod string, ledgerService LedgerServiceI) ResourceService {
	return ResourceService{
		didMethod:     didMethod,
		ledgerService: ledgerService,
	}
}

func (rds ResourceService) DereferenceResourceMetadata(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if err := rds.validateResourceRequest(did, &resourceId, contentType); err != nil {
		return nil, err
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}
	contentStream := types.NewDereferencedResourceList(did, []*resourceTypes.ResourceHeader{resource.Header})
	return &types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceCollectionResources(did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if err := rds.validateResourceRequest(did, nil, contentType); err != nil {
		return nil, err
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}
	contentStream := types.NewDereferencedResourceList(did, resources)
	return &types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceResourceData(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if err := rds.validateResourceRequest(did, &resourceId, contentType); err != nil {
		return nil, err
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}
	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}
	result := types.DereferencedResourceData(resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return &types.DidDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) validateResourceRequest(did string, resourceId *string, contentType types.ContentType) *types.IdentityError {
	if !contentType.IsSupported() {
		return types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}
	if !cheqdUtils.IsValidDID(did, rds.didMethod, rds.ledgerService.GetNamespaces()) || (resourceId != nil && !utils.IsValidResourceId(*resourceId)) {
		return types.NewInvalidDIDUrlError(did, contentType, nil, true)
	}
	return nil
}
