package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
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

	if utils.IsResourceHeaderPath(path) {
		if !dereferenceOptions.Accept.IsSupported() {
			dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
			return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
		}
		cotentStream, dereferenceMetadata = rds.dereferenceHeader(path, did, dereferenceOptions)
	} else if utils.IsCollectionResourcesPath(path) {
		if !dereferenceOptions.Accept.IsSupported() {
			dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
			return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
		}
		cotentStream, dereferenceMetadata = rds.dereferenceCollectionResources(did, dereferenceOptions)
	} else if utils.IsResourceDataPath(path) {
		cotentStream, dereferenceMetadata = rds.dereferenceResourceData(path, did, dereferenceOptions)
		if dereferenceOptions.Accept != dereferenceMetadata.ContentType {
			dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
			return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
		}
	} else {
		dereferenceMetadata = types.NewDereferencingMetadata(did, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
	}

	return types.DidDereferencing{ContentStream: cotentStream, DereferencingMetadata: dereferenceMetadata}
}

func (rds ResourceDereferenceService) dereferenceHeader(path string, did string, dereferenceOptions types.DereferencingOption) (*types.DereferencedResource, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")

	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return &types.DereferencedResource{}, dereferenceMetadata
	}
	return types.NewDereferencedResource(resource.Header), dereferenceMetadata
}

func (rds ResourceDereferenceService) dereferenceCollectionResources(did string, dereferenceOptions types.DereferencingOption) (*types.DereferencedResourceList, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resources, dereferencingError := rds.ledgerService.QueryCollectionResources(did)
	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return &types.DereferencedResourceList{}, dereferenceMetadata
	}
	return types.NewDereferencedResourceList(resources), dereferenceMetadata
}

func (rds ResourceDereferenceService) dereferenceResourceData(path string, did string, dereferenceOptions types.DereferencingOption) (*types.DereferencedResourceData, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return &types.DereferencedResourceData{}, dereferenceMetadata
	}
	result := types.DereferencedResourceData(resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return &result, dereferenceMetadata
}
