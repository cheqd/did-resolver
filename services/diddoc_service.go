package services

import (
	"net/url"
	"strings"

	didTypes "github.com/cheqd/cheqd-node/x/did/types"

	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog/log"
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

func (dds DIDDocService) ProcessDIDRequest(did string, fragmentId string, queries url.Values, flag *string, contentType types.ContentType) (types.ResolutionResultI, *types.IdentityError) {
	log.Trace().Msgf("ProcessDIDRequest %s, %s, %s", did, fragmentId, queries)
	var result types.ResolutionResultI
	var err *types.IdentityError
	var isDereferencing bool

	version := ""
	if len(queries) > 0 {
		version = queries.Get("versionId")
		if version == "" {
			return nil, types.NewRepresentationNotSupportedError(did, contentType, nil, true)
		}
	}

	switch {
	case flag != nil:
		return nil, types.NewRepresentationNotSupportedError(did, contentType, nil, true)
	case fragmentId != "":
		log.Trace().Msgf("Dereferencing %s, %s, %s", did, fragmentId, queries)
		result, err = dds.dereferenceSecondary(did, version, fragmentId, contentType)
		isDereferencing = true
	default:
		log.Trace().Msgf("Resolving %s", did)
		result, err = dds.Resolve(did, version, contentType)
		isDereferencing = false
	}

	if err != nil {
		err.IsDereferencing = isDereferencing
		return nil, err
	}
	return result, nil
}

func (dds DIDDocService) Resolve(did string, version string, contentType types.ContentType) (*types.DidResolution, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, false)
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

	protoDidDocWithMetadata, err := dds.ledgerService.QueryDIDDoc(did, version)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	resolvedMetadata, mErr := dds.resolveMetadata(did, *protoDidDocWithMetadata.Metadata, contentType)
	if mErr != nil {
		mErr.ContentType = contentType
		return nil, mErr
	}
	didDoc := types.NewDidDoc(*protoDidDocWithMetadata.DidDoc)
	result := types.DidResolution{Did: &didDoc, Metadata: *resolvedMetadata, ResolutionMetadata: didResolutionMetadata}
	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.AddContext(types.DIDSchemaJSONLD)
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
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, false)
	}

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

	contentStream := types.NewResolutionDidDocMetadata(did, *protoDidDocWithMetadata.Metadata, resources)

	return &types.ResourceDereferencing{Context: context, ContentStream: &contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (dds DIDDocService) GetAllDidDocVersionsMetadata(did string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	dereferenceMetadata := types.NewDereferencingMetadata(did, contentType, "")

	versions, err := dds.ledgerService.QueryAllDidDocVersionsMetadata(did)
	if err != nil {
		return nil, err
	}

	if len(versions) == 0 {
		return nil, types.NewNotFoundError(did, contentType, err, false)
	}

	var context string
	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		context = types.ResolutionSchemaJSONLD
	}

	contentStream := types.NewDereferencedDidVersionsList(versions)

	return &types.DidDereferencing{Context: context, ContentStream: contentStream, DereferencingMetadata: dereferenceMetadata}, nil
}

func (dds DIDDocService) dereferenceSecondary(did string, version string, fragmentId string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	didResolution, err := dds.Resolve(did, version, contentType)
	if err != nil {
		err.IsDereferencing = true
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

func (dds DIDDocService) resolveMetadata(did string, metadata didTypes.Metadata, contentType types.ContentType) (*types.ResolutionDidDocMetadata, *types.IdentityError) {
	resources, err := dds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return nil, err
	}
	resolvedMetadata := types.NewResolutionDidDocMetadata(did, metadata, resources)
	return &resolvedMetadata, nil
}
