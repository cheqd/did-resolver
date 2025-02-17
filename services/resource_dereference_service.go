package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	"strings"

	"github.com/cheqd/did-resolver/types"
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

func (rds ResourceService) DereferenceResourceMetadata(did string, resourceId string, contentType types.ContentType) (*types.ResourceDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	metadata := types.NewDereferencedResource(did, resource.Metadata)
	return &types.ResourceDereferencing{Context: context, Metadata: metadata, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceCollectionResources(did string, contentType types.ContentType) (*types.ResourceDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	didDoc, err := rds.ledgerService.QueryDIDDoc(did, "")
	if err != nil {
		return nil, err
	}

	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	contentStream := types.NewResolutionDidDocMetadata(did, didDoc.Metadata, resources)

	return &types.ResourceDereferencing{Context: context, ContentStream: &contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) ResolveCollectionResources(did string, contentType types.ContentType) (*types.DidResolution, *types.IdentityError) {
	resolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

	didDoc, err := rds.ledgerService.QueryDIDDoc(did, "")
	if err != nil {
		return nil, err
	}

	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	metadata := types.NewResolutionDidDocMetadata(did, didDoc.Metadata, resources)

	return &types.DidResolution{Context: context, Metadata: metadata, ResolutionMetadata: resolutionMetadata}, nil
}

func (rds ResourceService) DereferenceResourceData(did string, resourceId string, contentType types.ContentType) (*types.ResourceDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	result := types.DereferencedResourceData(resource.Resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Metadata.MediaType)

	return &types.ResourceDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceResourceDataWithMetadata(did string, resourceId string, contentType types.ContentType) (*types.ResourceDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	var result types.ContentStreamI
	result = types.NewDereferencedResourceData(resource.Resource.Data)
	metadata := types.NewDereferencedResource(did, resource.Metadata)
	if dereferenceMetadata.ContentType == types.JSON || dereferenceMetadata.ContentType == types.TEXT {
		if res, err := types.NewResourceData(resource.Resource.Data); err == nil {
			result = res
		}
	}

	return &types.ResourceDereferencing{Context: context, ContentStream: result, Metadata: metadata, DereferencingMetadata: dereferenceMetadata}, nil
}
