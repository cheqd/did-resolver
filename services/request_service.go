package services

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	cheqdTypes "github.com/cheqd/cheqd-node/x/cheqd/types"
	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
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

func (rs RequestService) ProcessDIDRequest(did string, fragmentId string, queries url.Values, flag *string, contentType types.ContentType) (types.ResolutionResultI, *types.IdentityError) {
	log.Trace().Msgf("ProcessDIDRequest %s, %s, %s", did, fragmentId, queries)
	var result types.ResolutionResultI
	var err *types.IdentityError
	var isDereferencing bool
	if len(queries) > 0 || flag != nil {
		return nil, types.NewRepresentationNotSupportedError(did, contentType, nil, false)
	} else if fragmentId != "" {
		log.Trace().Msgf("Dereferencing %s, %s, %s", did, fragmentId, queries)
		result, err = rs.dereferenceSecondary(did, fragmentId, contentType)
		isDereferencing = true
	} else {
		log.Trace().Msgf("Resolving %s", did)
		result, err = rs.Resolve(did, contentType)
		isDereferencing = false
	}

	if err != nil {
		err.DefineDisplaying(isDereferencing)
		return nil, err
	}
	return result, nil
}

// https://w3c-ccg.github.io/did-resolution/#resolving
func (rs RequestService) Resolve(did string, contentType types.ContentType) (*types.DidResolution, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, false)
	}
	didResolutionMetadata := types.NewResolutionMetadata(did, contentType, "")

	if didMethod, _, _, _ := cheqdUtils.TrySplitDID(did); didMethod != rs.didMethod {
		return nil, types.NewMethodNotSupportedError(did, contentType, nil, false)
	}
	if !cheqdUtils.IsValidDID(did, "", rs.ledgerService.GetNamespaces()) {
		return nil, types.NewInvalidDIDError(did, contentType, nil, false)
	}

	protoDidDoc, metadata, err := rs.ledgerService.QueryDIDDoc(did)
	if err != nil {
		err.ContentType = contentType
		return nil, err
	}

	resolvedMetadata, mErr := rs.ResolveMetadata(did, *metadata, contentType)
	if mErr != nil {
		mErr.ContentType = contentType
		return nil, mErr
	}
	didDoc := types.NewDidDoc(*protoDidDoc)
	if didResolutionMetadata.ContentType == types.DIDJSONLD || didResolutionMetadata.ContentType == types.JSONLD {
		didDoc.AddContext(types.DIDSchemaJSONLD)
	} else {
		didDoc.RemoveContext()
	}
	return &types.DidResolution{Did: didDoc, Metadata: *resolvedMetadata, ResolutionMetadata: didResolutionMetadata}, nil
}

// https://w3c-ccg.github.io/did-resolution/#dereferencing
func (rs RequestService) dereferenceSecondary(did string, fragmentId string, contentType types.ContentType) (*types.DidDereferencing, *types.IdentityError) {
	if !contentType.IsSupported() {
		return nil, types.NewRepresentationNotSupportedError(did, types.JSON, nil, true)
	}

	didResolution, err := rs.Resolve(did, contentType)
	if err != nil {
		return nil, err
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
		return nil, types.NewNotFoundError(did, contentType, nil, true)
	}

	if contentType == types.DIDJSONLD || contentType == types.JSONLD {
		contentStream.AddContext(types.DIDSchemaJSONLD)
	} else {
		contentStream.RemoveContext()
	}

	return &types.DidDereferencing{
		ContentStream:         contentStream,
		Metadata:              metadata,
		DereferencingMetadata: types.DereferencingMetadata(didResolution.ResolutionMetadata),
	}, nil
}

func (rs RequestService) ResolveMetadata(did string, metadata cheqdTypes.Metadata, contentType types.ContentType) (*types.ResolutionDidDocMetadata, *types.IdentityError) {
	if metadata.Resources == nil {
		return types.NewResolutionDidDocMetadata(did, metadata, nil), nil
	}
	resources, err := rs.ledgerService.QueryCollectionResources(did)
	if err != nil {
		return nil, err
	}
	return types.NewResolutionDidDocMetadata(did, metadata, resources), nil
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
		return err
	}

	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result, rErr := rs.ProcessDIDRequest(did, fragmentId, queries, flag, requestedContentType)
	if rErr != nil {
		return rErr
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(http.StatusOK, result, "  ")
}

func (rs RequestService) DereferenceResourceMetadata(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result, err := rs.resourceDereferenceService.DereferenceHeader(resourceId, did, requestedContentType)
	if err != nil {
		err.DefineDisplaying(true)
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(http.StatusOK, result, "  ")
}

func (rs RequestService) DereferenceResourceData(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result, err := rs.resourceDereferenceService.DereferenceResourceData(resourceId, did, requestedContentType)
	if err != nil {
		err.DefineDisplaying(true)
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.Blob(http.StatusOK, result.GetContentType(), result.GetBytes())
}

func (rs RequestService) DereferenceCollectionResources(c echo.Context) error {
	did := c.Param("did")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	resolutionResponse, err := rs.resourceDereferenceService.DereferenceCollectionResources(did, requestedContentType)
	if err != nil {
		err.DefineDisplaying(true)
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, resolutionResponse.GetContentType())
	return c.JSONPretty(http.StatusOK, resolutionResponse, "  ")
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
