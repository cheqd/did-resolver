package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
)

type ResourceDereferenceService struct {
	ledgerService LedgerServiceI
	didDocService DIDDocService
}

func NewResourceDereferenceService(ledgerService LedgerServiceI, didDocService DIDDocService) ResourceDereferenceService {
	return ResourceDereferenceService{
		ledgerService: ledgerService,
		didDocService: didDocService,
	}
}

func (rds ResourceDereferenceService) DereferenceHeader(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resource, err := rds.ledgerService.QueryResource(did, resourceId)
	if err != nil {
		return nil, err
	}
	contentStream := types.NewDereferencedResourceList(did, []*resourceTypes.ResourceHeader{resource.Header})
	return &types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceDereferenceService) DereferenceCollectionResources(did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return nil, err
	}
	contentStream := types.NewDereferencedResourceList(did, resources)
	return &types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceDereferenceService) DereferenceResourceData(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resource, err := rds.ledgerService.QueryResource(did, resourceId)
	if err != nil {
		return nil, err
	}
	result := types.DereferencedResourceData(resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return &types.DidDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}, nil
}
