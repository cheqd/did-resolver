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


func (rs RequestService) dereferenceService(didUrl string, dereferenceOptions types.DereferencingOption) (types.DidDereferencing, error) {
	did, _, query, fragmentId, _ := cheqdUtils.TrySplitDIDUrl(didUrl)
	didResolution := rs.Resolve(did, types.ResolutionOption(dereferenceOptions))

	dereferencingMetadata := types.DereferencingMetadata(didResolution.ResolutionMetadata)
	if dereferencingMetadata.ResolutionError != "" {
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	queryUrl, err := url.Parse("?" + query)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	parseQuery:= queryUrl.Query()

	queryId := parseQuery.Get("service")
	if queryId == "" {
		dereferencingMetadata = types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.RepresentationNotSupportedError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	service := rs.didDocService.GetDIDService(queryId, didResolution.Did)
	if service == nil {
		dereferencingMetadata = types.NewDereferencingMetadata(didUrl, dereferenceOptions.Accept, types.NotFoundError)
		return types.DidDereferencing{DereferencingMetadata: dereferencingMetadata}, nil
	}

	serviceEndpoint := CreateServiceEndpoint(parseQuery.Get("relativeRef"), fragmentId, service.ServiceEndpoint)
	metadata := types.TransformToFragmentMetadata(didResolution.Metadata)

	jsonFragment, err := json.Marshal(serviceEndpoint)
	if err != nil {
		return types.DidDereferencing{}, err
	}
	contentStream := json.RawMessage(jsonFragment)

	dereferencingMetadata = types.NewDereferencingMetadata(did, dereferenceOptions.Accept, "")
	return types.DidDereferencing{ContentStream: contentStream, Metadata: metadata, DereferencingMetadata: dereferencingMetadata}, nil
}
