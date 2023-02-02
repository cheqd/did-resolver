package services

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	didUtils "github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
)

type RequestService struct {
	didMethod                  string
	ledgerService              LedgerServiceI
	didDocService              DIDDocService
	resourceDereferenceService ResourceService
}

func NewRequestService(didMethod string, ledgerService LedgerServiceI) RequestService {
	return RequestService{
		didMethod:                  didMethod,
		ledgerService:              ledgerService,
		didDocService:              NewDIDDocService(didMethod, ledgerService),
		resourceDereferenceService: NewResourceService(didMethod, ledgerService),
	}
}

func (rs RequestService) ResolveDIDDoc(c echo.Context) error {
	splitDID := strings.Split(c.Param("did"), "#")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))

	did, err := url.QueryUnescape(splitDID[0])
	if err != nil {
		return types.NewInvalidDIDUrlError(splitDID[0], requestedContentType, err, true)
	}

	var fragmentId string
	if len(splitDID) == 2 {
		fragmentId = splitDID[1]
	}

	queryRaw, flag := prepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	//nolint: nestif
	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did

			if fragmentId != "" {
				path += url.PathEscape(fmt.Sprintf("#%s", fragmentId))
			}

			if queryRaw != "" {
				path += fmt.Sprintf("?%s", queryRaw)
			}

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	result, rErr := rs.didDocService.ProcessDIDRequest(did, fragmentId, queries, flag, requestedContentType)
	if rErr != nil {
		return rErr
	}

	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func (rs RequestService) ResolveDIDDocVersion(c echo.Context) error {
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))

	did, err := getDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), requestedContentType, err, true)
	}

	version := c.Param("version")

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did + types.DID_VERSION_PATH + version

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	result, rErr := rs.didDocService.Resolve(did, version, requestedContentType)
	if rErr != nil {
		return rErr
	}

	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func (rs RequestService) ResolveAllDidDocVersionsMetadata(c echo.Context) error {
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err := getDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), requestedContentType, err, true)
	}

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did + types.DID_VERSIONS_PATH

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	result, rErr := rs.didDocService.GetAllDidDocVersionsMetadata(did, requestedContentType)
	if rErr != nil {
		return rErr
	}

	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())

	return c.JSONPretty(http.StatusOK, result, "  ")
}

// DereferenceResourceMetadata godoc
//
//	@Summary		Fetch Resource-specific metadata
//	@Description	Get metadata for a specific Resource within a DID Resource Collection
//	@Tags			Resource Resolution
//	@Accept			application/did+ld+json,application/ld+json,application/did+json
//	@Produce		application/did+ld+json,application/ld+json,application/did+json
//	@Param			did			path		string	true	"Full DID with unique identifier"
//	@Param			resourceId	path		string	true	"Resource-specific unique identifier"
//	@Success		200			{object}	types.DidDereferencing
//	@Failure		400			{object}	types.IdentityError
//	@Failure		404			{object}	types.IdentityError
//	@Failure		406			{object}	types.IdentityError
//	@Failure		500			{object}	types.IdentityError
//	@Router			/{did}/resources/{resourceId}/metadata [get]
func (rs RequestService) DereferenceResourceMetadata(c echo.Context) error {
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err := getDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), requestedContentType, err, true)
	}
	resourceId := c.Param("resource")

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did + types.RESOURCE_PATH + resourceId + "/metadata"

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	result, errI := rs.resourceDereferenceService.DereferenceResourceMetadata(resourceId, did, requestedContentType)
	if errI != nil {
		errI.IsDereferencing = true
		return errI
	}

	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())

	return c.JSONPretty(http.StatusOK, result, "  ")
}

func (rs RequestService) DereferenceResourceData(c echo.Context) error {
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err := getDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), requestedContentType, err, true)
	}
	resourceId := c.Param("resource")

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did + types.RESOURCE_PATH + resourceId

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	result, errI := rs.resourceDereferenceService.DereferenceResourceData(resourceId, did, requestedContentType)
	if errI != nil {
		errI.IsDereferencing = true
		return errI
	}

	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())

	return c.Blob(http.StatusOK, result.GetContentType(), result.GetBytes())
}

func (rs RequestService) DereferenceCollectionResources(c echo.Context) error {
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	did, err := getDidParam(c)
	if err != nil {
		return types.NewInvalidDIDUrlError(c.Param("did"), requestedContentType, err, true)
	}

	didMethod, _, identifier, _ := didUtils.TrySplitDID(did)
	if didMethod != rs.didDocService.didMethod {
		return types.NewMethodNotSupportedError(did, requestedContentType, nil, false)
	}

	if !didUtils.IsValidDID(did, "", rs.didDocService.ledgerService.GetNamespaces()) {
		err := didUtils.ValidateDID(did, "", rs.didDocService.ledgerService.GetNamespaces())
		if err.Error() == types.NewInvalidIdentifierError().Error() && utils.IsMigrationNeeded(identifier) {
			did = utils.MigrateDID(did)
			path := types.RESOLVER_PATH + did + types.DID_METADATA

			return c.Redirect(http.StatusMovedPermanently, path)
		} else {
			return types.NewInvalidDIDError(did, requestedContentType, nil, false)
		}
	}

	resolutionResponse, errI := rs.resourceDereferenceService.DereferenceCollectionResources(did, requestedContentType)
	if errI != nil {
		errI.IsDereferencing = true
		return errI
	}

	c.Response().Header().Set(echo.HeaderContentType, resolutionResponse.GetContentType())

	return c.JSONPretty(http.StatusOK, resolutionResponse, "  ")
}

func getContentType(accept string) types.ContentType {
	typeList := strings.Split(accept, ",")
	for _, cType := range typeList {
		result := types.ContentType(strings.Split(cType, ";")[0])
		if result == "*/*" || result == types.JSONLD {
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

func getDidParam(c echo.Context) (string, error) {
	return url.QueryUnescape(c.Param("did"))
}
