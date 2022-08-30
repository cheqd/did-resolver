package services

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
	resourceTypes "github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cheqd/did-resolver/types"
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

func (rs RequestService) ProcessDIDRequest(did string, fragmentId string, queries url.Values, flag *string, contentType types.ContentType) types.ResolutionResultI {
	log.Trace().Msgf("ProcessDIDRequest %s, %s, %s", did, fragmentId, queries)
	var result types.ResolutionResultI
	if len(queries) > 0 || flag != nil {
		dereferencingMetadata := types.NewDereferencingMetadata(did, contentType, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	} else if fragmentId != "" {
		log.Trace().Msgf("Dereferencing %s, %s, %s", did, fragmentId, queries)
		result = rs.Dereference(did, fragmentId, contentType)
	} else {
		log.Trace().Msgf("Resolving %s", did)
		result = rs.Resolve(did, contentType)
	}
	return result
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, contentType types.ContentType) types.DidResolution {
	if !contentType.IsSupported() {
		return types.DidResolution{ResolutionMetadata: types.NewResolutionMetadata(did, types.JSON, types.RepresentationNotSupportedError)}
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

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
func (rs RequestService) Dereference(did string, fragmentId string, contentType types.ContentType) types.DidDereferencing {

	didDereferencing := rs.dereferenceSecondary(did, fragmentId, contentType)

	if didDereferencing.DereferencingMetadata.ResolutionError != "" {
		didDereferencing.ContentStream = nil
		didDereferencing.Metadata = types.ResolutionDidDocMetadata{}
		return didDereferencing
	}

	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		didDereferencing.ContentStream.AddContext(types.DIDSchemaJSONLD)
	} else {
		didDereferencing.ContentStream.RemoveContext()
	}

	return didDereferencing
}

func (rs RequestService) dereferencePrimary(path string, did string, contentType types.ContentType) types.DidDereferencing {
	// Only resource are available for primary dereferencing
	// return rs.resourceDereferenceService.DereferenceResource(path, did, contentType)
	return types.DidDereferencing{}
}

func (rs RequestService) dereferenceSecondary(did string, fragmentId string, contentType types.ContentType) types.DidDereferencing {
	if !contentType.IsSupported() {
		dereferencingMetadata := types.NewDereferencingMetadata(did, types.JSON, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}
	}

	didResolution := rs.Resolve(did, contentType)

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
		dereferencingMetadata := types.NewDereferencingMetadata(did, contentType, types.NotFoundError)
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

func (rs RequestService) ResolveDIDDoc(c echo.Context) error {
	splitedDID := strings.Split(c.Param("did"), "#")
	did := splitedDID[0]
	var fragmentId string
	if len(splitedDID) == 2 {
		fragmentId = splitedDID[1]
	}

	queryRaw, flag := prepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return errors.New(string(types.InternalError))
	}

	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result := rs.ProcessDIDRequest(did, fragmentId, queries, flag, requestedContentType)
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(result.GetStatus(), result, "  ")
}

func (rs RequestService) DereferenceResourceMetadata(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result := rs.resourceDereferenceService.DereferenceHeader(resourceId, did, requestedContentType)
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(result.GetStatus(), result, "  ")
}

func (rs RequestService) DereferenceResourceData(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result := rs.resourceDereferenceService.DereferenceResourceData(resourceId, did, requestedContentType)
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	if result.GetStatus() == http.StatusOK {
		return c.Blob(result.GetStatus(), result.GetContentType(), result.GetBytes())
	}
	return c.JSONPretty(result.GetStatus(), result, "  ")
}

func (rs RequestService) DereferenceCollectionResources(c echo.Context) error {
	did := c.Param("did")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	resolutionResponse := rs.resourceDereferenceService.DereferenceCollectionResources(did, requestedContentType)
	c.Response().Header().Set(echo.HeaderContentType, resolutionResponse.GetContentType())
	return c.JSONPretty(resolutionResponse.GetStatus(), resolutionResponse, "  ")
}

func getContentType(accept string) types.ContentType {
	tmp := strings.Split(accept, ",")
	for _, cType := range tmp {
		result := types.ContentType(strings.Split(cType, ";")[0])
		if result == "*/*" {
			result = types.DIDJSONLD
		}
		if result.IsSupported() {
			return result
		}
	}
	return ""
}

func prepareQueries(c echo.Context) (rawQuery string, flag *string) {
	rawQuery = c.Request().URL.RawQuery
	flagIndex := strings.LastIndex(rawQuery, "%23")
	if flagIndex == -1 || strings.Contains(rawQuery[flagIndex:], "&") {
		return rawQuery, nil
	}
	queryFlag := rawQuery[flagIndex:]
	return rawQuery[0:flagIndex], &queryFlag
}