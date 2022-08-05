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

func (rds ResourceDereferenceService) DereferenceResource(path string, did string, didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, int) {
	var jsonCotentStream string
	var dereferenceMetadata types.DereferencingMetadata

	if utils.IsResourceHeaderPath(path) {
		jsonCotentStream, dereferenceMetadata = rds.dereferenceHeader(path, did, dereferenceOptions)
	} else if utils.IsCollectionResourcesPath(path) || utils.IsCollectionResourcesPathRedirect(path) {
		jsonCotentStream, dereferenceMetadata = rds.dereferenceCollectionResources(path, did, dereferenceOptions)
	} else {
		dereferenceMetadata.ResolutionError = types.RepresentationNotSupportedError
	}

	var statusCode int
	if utils.IsCollectionResourcesPathRedirect(path) {
		statusCode = 301 // redirect code
	} else {
		statusCode = dereferenceMetadata.ResolutionError.GetStatusCode()
	}

	return types.DidDereferencing{ContentStream: json.RawMessage(jsonCotentStream), DereferencingMetadata: dereferenceMetadata}, statusCode
}

func (rds ResourceDereferenceService) dereferenceHeader(path string, did string, dereferenceOptions types.DereferencingOption) (string, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")

	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return "", dereferenceMetadata
	}
	var err error
	cotentStream, err := rds.didDocService.MarshallContentStream(resource.Header, dereferenceOptions.Accept)
	if err != nil {
		dereferenceMetadata.ResolutionError = types.InternalError
	}
	return cotentStream, dereferenceMetadata
}

func (rds ResourceDereferenceService) dereferenceCollectionResources(path string, did string, dereferenceOptions types.DereferencingOption) (string, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resources, dereferencingError := rds.ledgerService.QueryCollectionResources(did)
	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return "", dereferenceMetadata
	}
	jsonResources := []json.RawMessage{}
	for _, r := range resources {
		jsonR, err := rds.didDocService.MarshallContentStream(r, dereferenceOptions.Accept)
		if err != nil {
			dereferenceMetadata.ResolutionError = types.InternalError
			return "", dereferenceMetadata
		}
		jsonResources = append(jsonResources, json.RawMessage(jsonR))
	}
	cotentStream, err := json.Marshal(jsonResources)
	if err != nil {
		dereferenceMetadata.ResolutionError = types.InternalError
		return "", dereferenceMetadata
	}
	return string(cotentStream), dereferenceMetadata
}

func (rds ResourceDereferenceService) DereferenceResourceData(path string, did string, didUrl string, dereferenceOptions types.DereferencingOption) ([]byte, types.DereferencingMetadata) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	resourceId := utils.GetResourceId(path)

	resource, dereferencingError := rds.ledgerService.QueryResource(did, resourceId)

	if dereferencingError != "" {
		dereferenceMetadata.ResolutionError = dereferencingError
		return []byte{}, dereferenceMetadata
	}
	dereferenceMetadata.ContentType = types.ContentType(resource.Header.MediaType)
	return resource.Data, dereferenceMetadata
}
