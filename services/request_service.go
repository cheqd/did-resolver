package services

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/types"
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
	result, rErr := rs.didDocService.ProcessDIDRequest(did, fragmentId, queries, flag, requestedContentType)
	if rErr != nil {
		return rErr
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(http.StatusOK, result, "  ")
}

// DereferenceResourceMetadata godoc
// @Summary      Resource metadata
// @Description  Get resource metadata without value by DID Doc
// @Tags         Dereferencing
// @Produce      json
// @Param        did path string true "Resource collection id. DID Doc Id" example("did:cheqd:testnet:MjYxNzYKMjYxNzYK")
// @Param        resourceId path string true "DID Resource identifier" example("60ad67be-b65b-40b8-b2f4-3923141ef382")
// @Success      200  {object}  types.DidDereferencing
// @Failure      400  {object}  types.DidDereferencing
// @Failure      404  {object}  types.DidDereferencing
// @Failure      406  {object}  types.DidDereferencing
// @Failure      500  {object}  types.DidDereferencing
// @Router       /1.0/identifiers/{did}/resources/{resourceId}/metadata [get]
func (rs RequestService) DereferenceResourceMetadata(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result, err := rs.resourceDereferenceService.DereferenceResourceMetadata(resourceId, did, requestedContentType)
	if err != nil {
		err.IsDereferencing = true
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.JSONPretty(http.StatusOK, result, "  ")
}

// DereferenceResourceData godoc
// @Summary      Resource value
// @Description  Get resource value without dereferencing wrappers
// @Tags         Dereferencing
// @Produce      */*
// @Param        did path string true "Resource collection id. DID Doc Id" example("did:cheqd:testnet:MjYxNzYKMjYxNzYK")
// @Param        resourceId path string true "DID Resource identifier" example("60ad67be-b65b-40b8-b2f4-3923141ef382")
// @Success      200  {object}  []byte
// @Failure      400  {object}  types.DidDereferencing
// @Failure      404  {object}  types.DidDereferencing
// @Failure      406  {object}  types.DidDereferencing
// @Failure      500  {object}  types.DidDereferencing 
// @Router       /1.0/identifiers/{did}/resources/{resourceId} [get]
func (rs RequestService) DereferenceResourceData(c echo.Context) error {
	did := c.Param("did")
	resourceId := c.Param("resource")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	result, err := rs.resourceDereferenceService.DereferenceResourceData(resourceId, did, requestedContentType)
	if err != nil {
		err.IsDereferencing = true
		return err
	}
	c.Response().Header().Set(echo.HeaderContentType, result.GetContentType())
	return c.Blob(http.StatusOK, result.GetContentType(), result.GetBytes())
}

// DereferenceCollectionResources godoc
// @Summary      Collection resources
// @Description  Get a list of all collection resources metadata
// @Tags         Dereferencing
// @Produce      json
// @Param        did path string true "Resource collection id. DID Doc Id" example("did:cheqd:testnet:MjYxNzYKMjYxNzYK") validate(optional)
// @Success      200  {object}  types.DidDereferencing
// @Failure      400  {object}  types.DidDereferencing
// @Failure      404  {object}  types.DidDereferencing
// @Failure      406  {object}  types.DidDereferencing
// @Failure      500  {object}  types.DidDereferencing
// @Router       /1.0/identifiers/{did}/resources/all [get]
func (rs RequestService) DereferenceCollectionResources(c echo.Context) error {
	did := c.Param("did")
	requestedContentType := getContentType(c.Request().Header.Get(echo.HeaderAccept))
	resolutionResponse, err := rs.resourceDereferenceService.DereferenceCollectionResources(did, requestedContentType)
	if err != nil {
		err.IsDereferencing = true
		return err
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
