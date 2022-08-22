package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
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

func (rds ResourceDereferenceService) DereferenceResource(path string, did string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	var cotentStream types.ContentStreamI
	var dereferenceMetadata types.DereferencingMetadata

	if !dereferenceOptions.Accept.IsSupported() && !utils.IsResourceDataPath(path) {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	if utils.IsResourceHeaderPath(path) {
		cotentStream, dereferenceMetadata = rds.DereferenceHeader(path, did, dereferenceOptions)
	} else if utils.IsCollectionResourcesPath(path) {
		cotentStream, dereferenceMetadata = rds.DereferenceCollectionResources(did, dereferenceOptions)
	} else if utils.IsResourceDataPath(path) {
		cotentStream, dereferenceMetadata = rds.DereferenceResourceData(path, did, dereferenceOptions)
	} else {
		dereferenceMetadata = types.NewDereferencingMetadata(did, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
	}

	return types.DidDereferencing{ContentStream: cotentStream, DereferencingMetadata: dereferenceMetadata}
}

func (rds ResourceDereferenceService) DereferenceHeader(resourceId string, did string, contentType types.ContentType) types.DidDereferencing {
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
