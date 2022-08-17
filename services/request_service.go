package services

import (
	"github.com/rs/zerolog/log"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type RequestService struct {
	didMethod                  string
	ledgerService              LedgerServiceI
	didDocService              DIDDocService
	resourceDereferenceService ResourceDereferenceService
}

func NewRequestService(didMethod string, ledgerService LedgerServiceI) RequestService {
	didDocService := DIDDocService{}
	return RequestService{
		didMethod:                  didMethod,
		ledgerService:              ledgerService,
		didDocService:              didDocService,
		resourceDereferenceService: NewResourceDereferenceService(ledgerService, didDocService),
	}
}

func (rs RequestService) ProcessDIDRequest(didUrl string, resolutionOptions types.ResolutionOption) types.ResolutinResultI {
	var result types.ResolutinResultI
	did, path, query, fragmentId, _ := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Warn().Msgf("Query %s %s %s %s ", did, path, query, fragmentId)
	if utils.IsDidUrl(didUrl) {
		log.Trace().Msgf("Dereferencing %s", didUrl)
		result = rs.Dereference(didUrl, types.DereferencingOption(resolutionOptions))
	} else {
		log.Trace().Msgf("Resolving %s", didUrl)
		result = rs.Resolve(didUrl, resolutionOptions)
	}
	return result
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, resolutionOptions types.ResolutionOption) types.DidResolution {
	if !resolutionOptions.Accept.IsSupported() {
		return types.DidResolution{ResolutionMetadata: types.NewResolutionMetadata(did, types.JSON, types.RepresentationNotSupportedError)}
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, resolutionOptions.Accept, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != rs.didMethod {
		didResolutionMetadata.ResolutionError = types.MethodNotSupportedError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}
	if !cheqdUtils.IsValidDID(did, "", rs.ledgerService.GetNamespaces()) {
		didResolutionMetadata.ResolutionError = types.InvalidDIDError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}

	}

	protoDidDoc, metadata, isFound, err := rs.ledgerService.QueryDIDDoc(did)
	if err != nil {
		didResolutionMetadata.ResolutionError = types.InternalError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}

	resolvedMetadata, errorType := rs.ResolveMetadata(did, metadata)
	if errorType != "" {
		didResolutionMetadata.ResolutionError = types.InternalError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}

	if !isFound {
		didResolutionMetadata.ResolutionError = types.NotFoundError
		return types.DidResolution{ResolutionMetadata: didResolutionMetadata}
	}
	didDoc := types.NewDidDoc(protoDidDoc)
	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.AddContext(types.DIDSchemaJSONLD)
	} else {
		didDoc.RemoveContext()
	}
	return types.DidResolution{Did: didDoc, Metadata: resolvedMetadata, ResolutionMetadata: didResolutionMetadata}
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (rs RequestService) Dereference(didUrl string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	did, path, query, fragmentId, err := cheqdUtils.TrySplitDIDUrl(didUrl)
	log.Info().Msgf("did: %s, path: %s, query: %s, fragmentId: %s", did, path, query, fragmentId)

	if err != nil || !cheqdUtils.IsValidDIDUrl(didUrl, "", []string{}) {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.InvalidDIDUrlError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	// TODO: implement
	if query != "" {
		dereferencingMetadata := types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	var didDereferencing types.DidDereferencing
	if path != "" {
		didDereferencing = rs.dereferencePrimary(path, did, dereferenceOptions)
	} else {
		didDereferencing = rs.dereferenceSecondary(did, fragmentId, dereferenceOptions)
	}

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		didDereferencing.ContentStream = nil
		didDereferencing.Metadata = types.ResolutionDidDocMetadata{}
		return didDereferencing
	}

	if dereferenceOptions.Accept == types.DIDJSONLD || dereferenceOptions.Accept == types.JSONLD {

		didDereferencing.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else {
		didDereferencing.ContentStream.RemoveContext()
	}

	return didDereferencing
}

func (rs RequestService) dereferencePrimary(path string, did string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	// Only resource are available for primary dereferencing
	return rs.resourceDereferenceService.DereferenceResource(path, did, dereferenceOptions)
}

func (rs RequestService) dereferenceSecondary(did string, fragmentId string, dereferenceOptions types.DereferencingOption) types.DidDereferencing {
	if !dereferenceOptions.Accept.IsSupported() {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	didResolution := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))

	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	metadata := didResolution.Metadata

	var contentStream types.ContentStreamI
	if fragmentId != "" {
		contentStream = rs.didDocService.GetDIDFragment(fragmentId, *didResolution.Did)
		metadata = types.TransformToFragmentMetadata(metadata)
	} else {
		contentStream = didResolution.Did
	}

	if contentStream == nil {
		dereferencingMetadata := types.NewDereferencingMetadata(did, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}
	return types.DidDereferencing{ContentStream: contentStream, Metadata: metadata, DereferencingMetadata: dereferencingMetadata}
}

func (rs RequestService) ResolveMetadata(did string, metadata cheqdTypes.Metadata) (types.ResolutionDidDocMetadata, types.ErrorType) {
	if metadata.Resources == nil {
		return types.NewResolutionDidDocMetadata(did, metadata, []*resourceTypes.ResourceHeader{}), ""
	}
	resources, errorType := rs.ledgerService.QueryCollectionResources(did)
	if errorType != "" {
		return types.ResolutionDidDocMetadata{}, errorType
	}
	return types.NewResolutionDidDocMetadata(did, metadata, resources), ""
}
