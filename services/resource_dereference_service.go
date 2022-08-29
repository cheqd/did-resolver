package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog/log"
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

func (rds ResourceDereferenceService) DereferenceHeader(resourceId string, did string, contentType types.ContentType) types.DidDereferencing {
	if !contentType.IsSupported() {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)
	log.Warn().Msgf("dereferencingError: %s", dereferencingError)
	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return types.DidDereferencing{DereferencingMetadata: dereferenceMetadata}
	}
	contentStream := types.NewDereferencedResourceList(did, []*resourceTypes.ResourceHeader{resource.Header})
	return types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}
}

func (rds ResourceDereferenceService) DereferenceCollectionResources(did string, contentType types.ContentType) types.DidDereferencing {
	if !contentType.IsSupported() {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")
	resources, dereferencingError := rds.ledgerService.QueryCollectionResources(did)
	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return types.DidDereferencing{DereferencingMetadata: dereferenceMetadata}
	}
	contentStream := types.NewDereferencedResourceList(did, resources)
	return types.DidDereferencing{ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}
}

func (rds ResourceDereferenceService) DereferenceResourceData(resourceId string, did string, contentType types.ContentType) types.DidDereferencing {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return types.DidDereferencing{DereferencingMetadata: dereferenceMetadata}
	}
	result := types.DereferencedResourceData(resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return types.DidDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}
}
