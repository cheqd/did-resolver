package services

import (
	// jsonpb Marshaller is deprecated, but is needed because there's only one way to proto
	// marshal in combination with our proto generator version

	"strings"

	migrations "github.com/cheqd/cheqd-node/app/migrations/helpers"
	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/google/uuid"
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

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rds.didMethod {
		return nil, types.NewMethodNotSupportedError(did, contentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rds.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rds.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsValidV1ID(identifier) {
			did = migrations.MigrateIndyStyleDid(did)
		} else {
			return nil, types.NewInvalidDIDError(did, contentType, nil, false)
		}
	}

	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		_, parsingerr := uuid.Parse(identifier)
		if parsingerr == nil {
			did = migrations.MigrateUUIDDid(did)
			resource, err = rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
			if err != nil {
				err.ContentType = contentType
				return nil, err
			}
		} else {
			err.ContentType = contentType
			return nil, err
		}
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

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rds.didMethod {
		return nil, types.NewMethodNotSupportedError(did, contentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rds.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rds.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsValidV1ID(identifier) {
			did = migrations.MigrateIndyStyleDid(did)
		} else {
			return nil, types.NewInvalidDIDError(did, contentType, nil, false)
		}
	}

	resources, err := rds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		_, parsingerr := uuid.Parse(identifier)
		if parsingerr == nil {
			did = migrations.MigrateUUIDDid(did)
			resources, err = rds.ledgerService.QueryCollectionResources(did)
			if err != nil {
				err.ContentType = contentType
				return nil, err
			}
		} else {
			err.ContentType = contentType
			return nil, err
		}
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

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rds.didMethod {
		return nil, types.NewMethodNotSupportedError(did, contentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rds.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rds.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsValidV1ID(identifier) {
			did = migrations.MigrateIndyStyleDid(did)
		} else {
			return nil, types.NewInvalidDIDError(did, contentType, nil, false)
		}
	}

	resource, err := rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
	if err != nil {
		_, parsingerr := uuid.Parse(identifier)
		if parsingerr == nil {
			did = migrations.MigrateUUIDDid(did)
			resource, err = rds.ledgerService.QueryResource(did, strings.ToLower(resourceId))
			if err != nil {
				err.ContentType = contentType
				return nil, err
			}
		} else {
			err.ContentType = contentType
			return nil, err
		}
	}

	result := types.DereferencedResourceData(resource.Resource.Data)
	dereferenceMetadata.ContentType = types.ContentType(resource.Metadata.MediaType)

	return &types.DidDereferencing{ContentStream: &result, DereferencingMetadata: dereferenceMetadata}, nil
}
