package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	"strings"

	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
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

func (rds ResourceService) DereferenceResourceMetadata(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}

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

	contentStream := types.NewDereferencedResourceList(did, []*resourceTypes.Metadata{resource.Metadata})

	return &types.DidDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceCollectionResources(did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}

	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	contentStream := types.NewDereferencedResourceList(did, resources)

	return &types.DidDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (rds ResourceService) DereferenceResourceData(resourceId string, did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
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

	result := types.DereferencedResourceData(resource.Resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Metadata.MediaType)

	return &types.DidDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}, nil
}
