package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version
	"encoding/json"

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
	var cotentStream []byte
	var dereferenceMetadata types.DereferencingMetadata

	if utils.IsResourceHeaderPath(path) {
		cotentStream, dereferenceMetadata = rds.dereferenceHeader(path, did, dereferenceOptions)
	} else if utils.IsCollectionResourcesPath(path) {
		cotentStream, dereferenceMetadata = rds.dereferenceCollectionResources(did, dereferenceOptions)
	} else if utils.IsResourceDataPath(path) {
		cotentStream, dereferenceMetadata = rds.dereferenceResourceData(path, did, dereferenceOptions)
	} else {
		dereferenceMetadata = types.NewDereferencingMetadata(did, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
	}

	return types.DidDereferencing{ContentStream: cotentStream, DereferencingMetadata: dereferenceMetadata}
}

func (rds ResourceDereferenceService) dereferenceHeader(path string, did string, dereferenceOptions types.DereferencingOption) ([]byte, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")

	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return []byte(nil), dereferenceMetadata
	}
	var err error
	cotentStream, err := rds.didDocService.MarshallContentStream(resource.Header, dereferenceOptions.Accept)
	if err != nil {
		dereferenceMetadata.ResolutionError = types.InternalError
	}
	return []byte(cotentStream), dereferenceMetadata
}

func (rds ResourceDereferenceService) dereferenceCollectionResources(did string, dereferenceOptions types.DereferencingOption) ([]byte, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resources, dereferencingError := rds.ledgerService.QueryCollectionResources(did)
	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return []byte(nil), dereferenceMetadata
	}
	jsonResources := []json.RawMessage{}
	for _, r := range resources {
		jsonR, err := rds.didDocService.MarshallContentStream(r, dereferenceOptions.Accept)
		if err != nil {
			dereferenceMetadata.ResolutionError = types.InternalError
			return []byte(nil), dereferenceMetadata
		}
		jsonResources = append(jsonResources, json.RawMessage(jsonR))
	}
	cotentStream, err := json.MarshalIndent(jsonResources, "", "  ")
	if err != nil {
		dereferenceMetadata.ResolutionError = types.InternalError
		return []byte(nil), dereferenceMetadata
	}
	return cotentStream, dereferenceMetadata
}

func (rds ResourceDereferenceService) dereferenceResourceData(path string, did string, dereferenceOptions types.DereferencingOption) ([]byte, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return []byte(nil), dereferenceMetadata
	}
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return resource.Data, dereferenceMetadata
}
