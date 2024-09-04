package services

import (
	"strings"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"

	"github.com/cheqd/did-resolver/types"
)

type DIDDocService struct {
	didMethod     string
	ledgerService LedgerServiceI
}

func NewDIDDocService(didMethod string, ledgerService LedgerServiceI) DIDDocService {
	return DIDDocService{
		didMethod:     didMethod,
		ledgerService: ledgerService,
	}
}

func (DIDDocService) GetDIDFragment(fragmentId string, didDoc types.DidDoc) types.ContentStreamI {
	for _, verMethod := range didDoc.VerificationMethod {
		if strings.Contains(verMethod.Id, fragmentId) {
			return &verMethod
		}
	}
	for _, service := range didDoc.Service {
		if strings.Contains(service.Id, fragmentId) {
			return &service
		}
	}

	return nil
}

func (dds DIDDocService) Resolve(did string, version string, contentType types.ContentType) (*types.DidResolution, *types.IdentityError) {
	didResolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

	protoDidDocWithMetadata, err := dds.ledgerService.QueryDIDDoc(did, version)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	resolvedMetadata, mErr := dds.resolveMetadata(did, protoDidDocWithMetadata.Metadata, contentType)
	if mErr != nil {
		mErr.ContentType = contentType
		return nil, mErr
	}
	didDoc := types.NewDidDoc(protoDidDocWithMetadata.DidDoc)
	result := types.DidResolution{Did: &didDoc, Metadata: *resolvedMetadata, ResolutionMetadata: didResolutionMetadata}
	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.AddContext(types.DIDSchemaJSONLD)

		if len(didDoc.Service) > 0 {
			didDoc.AddContext(types.DIFDIDConfigurationJSONLD)
		}

		for _, method := range didDoc.VerificationMethod {
			switch method.Type {
			case "Ed25519VerificationKey2020":
				didDoc.AddContext(types.Ed25519VerificationKey2020JSONLD)
			case "Ed25519VerificationKey2018":
				didDoc.AddContext(types.Ed25519VerificationKey2018JSONLD)
			case "JsonWebKey2020":
				didDoc.AddContext(types.JsonWebKey2020JSONLD)
			}
		}
		result.Context = types.ResolutionSchemaJSONLD
	} else {
		didDoc.RemoveContext()
	}

	return &result, nil
}

func (dds DIDDocService) GetDIDDocVersionsMetadata(did string, version string, contentType types.ContentType) (*types.ResourceDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	protoDidDocWithMetadata, err := dds.ledgerService.QueryDIDDoc(did, version)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	resources, err := dds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	contentStream := types.NewResolutionDidDocMetadata(did, protoDidDocWithMetadata.Metadata, resources)

	return &types.ResourceDereferencing{Context: context, ContentStream: &contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (dds DIDDocService) GetAllDidDocVersionsMetadata(did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	versions, err := dds.ledgerService.QueryAllDidDocVersionsMetadata(did)
	if err != nil {
		return nil, err
	}

	resources, err := dds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	if len(versions) == 0 {
		return nil, types.NewNotFoundError(did, contentType, err, false)
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	contentStream := types.NewDereferencedDidVersionsList(did, versions, resources)
	for i, version := range contentStream.Versions {
		filtered := contentStream.Versions.GetResourcesBeforeNextVersion(version.VersionId)
		contentStream.Versions[i].Resources = filtered
	}

	return &types.DidDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (dds DIDDocService) DereferenceSecondary(did string, version string, fragmentId string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	didResolution, err := dds.Resolve(did, version, contentType)
	if err != nil {
		return nil, err
	}

	metadata := didResolution.Metadata

	var contentStream types.ContentStreamI
	if fragmentId != "" {
		contentStream = dds.GetDIDFragment(fragmentId, *didResolution.Did)
		metadata = types.TransformToFragmentMetadata(metadata)
	} else {
		contentStream = didResolution.Did
	}

	if contentStream == nil {
		return nil, types.NewNotFoundError(did, contentType, nil, true)
	}

	result := types.DidDereferencing{
		ContentStream:         contentStream,
		Metadata:              metadata,
		DereferencingMetadata: types.DereferencingMetadata(didResolution.ResolutionMetadata),
	}

	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		contentStream.AddContext(types.DIDSchemaJSONLD)
		result.Context = types.ResolutionSchemaJSONLD
	} else {
		contentStream.RemoveContext()
	}

	return &result, nil
}

func (dds DIDDocService) resolveMetadata(did string, metadata *didTypes.Metadata, contentType types.ContentType) (*types.ResolutionDidDocMetadata, *types.IdentityError) {
	resources, err := dds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return nil, err
	}

	resolvedMetadata := types.NewResolutionDidDocMetadata(did, metadata, resources)

	return &resolvedMetadata, nil
}
