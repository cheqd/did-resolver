package services

import (
	"net/url"
	"strings"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
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

	if len(queries) > 0 || flag != nil {
		return nil, types.NewRepresentationNotSupportedError(did, contentType, nil, true)
	} else if fragmentId != "" {
		log.Trace().Msgf("Dereferencing %s, %s, %s", did, fragmentId, queries)
		result, err = dds.dereferenceSecondary(did, fragmentId, contentType)
		isDereferencing = true
	} else {
		log.Trace().Msgf("Resolving %s", did)
		result, err = dds.Resolve(did, contentType)
		isDereferencing = false
	}

	if err != nil {
		err.IsDereferencing = isDereferencing
		return nil, err
	}
	return result, nil
}

// ResolveDIDDoc godoc
// @Summary      Resolve or dereferencing DID Doc
// @Description  Get DID Doc or its fragment
// @Tags         Resolution
// @Accept       application/did+ld+json,application/ld+json,application/did+json
// @Produce      application/did+ld+json,application/ld+json,application/did+json
// @Param        did path string true "DID Doc Id" example(did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY)
// @Success      200  {object}  types.DidResolution
// @Failure      400  {object}  types.DidResolution
// @Failure      404  {object}  types.DidResolution
// @Failure      406  {object}  types.DidResolution
// @Failure      500  {object}  types.DidResolution
// @Router       /1.0/identifiers/{did} [get]
func (dds DIDDocService) Resolve(did string, contentType types.ContentType) (*types.DidResolution, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, false)
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != dds.didMethod {
		return nil, types.NewMethodNotSupportedError(did, contentType, nil, false)
	}
	if !cheqdUtils.IsValidDID(did, "", dds.ledgerService.GetNamespaces()) {
		return nil, types.NewInvalidDIDError(did, contentType, nil, false)
	}

	protoDidDoc, metadata, err := dds.ledgerService.QueryDIDDoc(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	resolvedMetadata, mErr := dds.resolveMetadata(did, *metadata, contentType)
	if mErr != nil {
		mErr.ContentType = contentType
		return nil, mErr
	}
	didDoc := types.NewDidDoc(*protoDidDoc)
	result := types.DidResolution{Did: &didDoc, Metadata: *resolvedMetadata, ResolutionMetadata: didResolutionMetadata}
	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.AddContext(types.DIDSchemaJSONLD)
		result.Context = types.ResolutionSchemaJSONLD
	} else {
		didDoc.RemoveContext()
	}
	return &result, nil
}

// ResolveDIDDocDereferencing godoc
// @Summary      Resolve or dereferencing DID Doc
// @Description  Get DID Doc or its fragment
// @Tags         Dereferencing
// @Accept       application/did+ld+json,application/ld+json,application/did+json
// @Produce      application/did+ld+json,application/ld+json,application/did+json
// @Param        did path string true "DID Doc Id" example(did:cheqd:mainnet:zF7rhDBfUt9d1gJPjx7s1JXfUY7oVWkY)
// @Param        fragmentId path string false "`#` + DID Doc Verification Method or Service identifier" example(#key1)
// @Param        service query string false "Service id" example("service1")
// @Success      200  {object}  types.DidDereferencing
// @Failure      400  {object}  types.DidDereferencing
// @Failure      404  {object}  types.DidDereferencing
// @Failure      406  {object}  types.DidDereferencing
// @Failure      500  {object}  types.DidDereferencing
// @Router       /1.0/identifiers/{did}{fragmentId} [get]
func (dds DIDDocService) dereferenceSecondary(did string, fragmentId string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	didResolution, err := dds.Resolve(did, contentType)
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

func (dds DIDDocService) resolveMetadata(did string, metadata cheqdTypes.Metadata, contentType types.ContentType) (*types.ResolutionDidDocMetadata, *types.IdentityError) {
	if metadata.Resources == nil {
		resolvedMetadata := types.NewResolutionDidDocMetadata(did, metadata, nil)
		return &resolvedMetadata, nil
	}
	resources, err := dds.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return nil, err
	}
	resolvedMetadata := types.NewResolutionDidDocMetadata(did, metadata, resources)
	return &resolvedMetadata, nil
}
