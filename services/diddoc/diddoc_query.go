package diddoc

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v4"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/diddoc/queries"
	diddocQueries "github.com/cheqd/did-resolver/services/diddoc/queries/diddoc"
	resourceQueries "github.com/cheqd/did-resolver/services/diddoc/queries/resources"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/rs/zerolog/log"
)

type QueryDIDDocRequestService struct {
	services.BaseRequestService
	Profile      string
	FirstHandler queries.BaseQueryHandlerI
}

func (dd *QueryDIDDocRequestService) Setup(c services.ResolverContext) error {
	return nil
}

// function to prepare the Query DIDDoc Request with specific conditions
func (dd *QueryDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	// if profile is W3IDDIDRES then dereferencing is false
	acceptHeader := c.Request().Header.Get(echo.HeaderAccept)
	contentType, profile := services.GetPriorityContentType(acceptHeader, dd.AreResourceQueriesPlaced(c))

	dd.Profile = profile
	dd.RequestedContentType = contentType

	if profile == types.W3IDDIDRES {
		dd.IsDereferencing = false
	} else {
		dd.IsDereferencing = true
	}

	// Register query handlers
	return dd.RegisterQueryHandlers(c)
}

// function to validate the Query DIDDoc Request with specific conditions
func (dd *QueryDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	_, err := url.QueryUnescape(dd.GetDid())
	if err != nil {
		log.Debug().Msg(err.Error())
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.RequestedContentType, err, dd.IsDereferencing)
	}

	diff := types.AllSupportedQueries.DiffWithUrlValues(dd.Queries)
	if len(diff) > 0 {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	if dd.AreQueryValuesEmpty(c) {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	versionId := dd.GetQueryParam(types.VersionId)
	versionTime := dd.GetQueryParam(types.VersionTime)
	transformKeys := types.TransformKeysType(dd.GetQueryParam(types.TransformKeys))
	service := dd.GetQueryParam(types.ServiceQ)
	relativeRef := dd.GetQueryParam(types.RelativeRef)
	resourceId := dd.GetQueryParam(types.ResourceId)
	resourceVersionTime := dd.GetQueryParam(types.ResourceVersionTime)
	metadata := dd.GetQueryParam(types.Metadata)
	resourceMetadata := dd.GetQueryParam(types.ResourceMetadata)
	if string(transformKeys) != "" && (!transformKeys.IsSupported() || !types.IsSupportedWithCombinationTransformKeysQuery(dd.Queries)) {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// relativeRef should be only with service parameter also
	if relativeRef != "" && service == "" {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// service query is permitted only for diddoc queries
	if service != "" && dd.AreResourceQueriesPlaced(c) {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// metadata query is permitted only for diddoc queries and for resource queries if resourceMetadata is placed
	if metadata != "" && (dd.AreResourceQueriesPlaced(c) && resourceMetadata == "") {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// value if metadata can be only true or false
	if metadata != "" && metadata != "true" && metadata != "false" {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// value if resourceMetadata can be only true or false
	if resourceMetadata != "" && resourceMetadata != "true" && resourceMetadata != "false" {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// if profile is W3IDDIDURL then metadata should be true
	if resourceMetadata == "false" && dd.Profile == types.W3IDDIDURL {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	// Validate time format
	if versionTime != "" {
		_, err := utils.ParseFromStringTimeToGoTime(versionTime)
		if err != nil {
			log.Debug().Msg(err.Error())
			return types.NewInvalidDidUrlError(dd.GetDid(), dd.GetContentType(), err, dd.IsDereferencing)
		}
	}

	// Validate time format
	if resourceVersionTime != "" {
		_, err := utils.ParseFromStringTimeToGoTime(resourceVersionTime)
		if err != nil {
			log.Debug().Msg(err.Error())
			return types.NewInvalidDidUrlError(dd.GetDid(), dd.GetContentType(), err, dd.IsDereferencing)
		}
	}

	// Validate that versionId is UUID
	if versionId != "" && !utils.IsValidUUID(versionId) {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	// Validate that resourceId is UUID
	if resourceId != "" && !utils.IsValidUUID(resourceId) {
		return types.NewInvalidDidUrlError(dd.GetDid(), dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	// If there is only 1 query parameter and it's resourceVersionTime,
	// then we need to return RepresentationNotSupported error
	if len(dd.Queries) == 1 && resourceVersionTime != "" {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}

	return nil
}

func (dd QueryDIDDocRequestService) AreResourceQueriesPlaced(c services.ResolverContext) bool {
	return len(types.ResourceSupportedQueries.IntersectWithUrlValues(dd.Queries)) > 0
}

func (dd QueryDIDDocRequestService) AreDidResolutionQueries(c services.ResolverContext) bool {
	return len(types.DidResolutionQueries.DiffWithUrlValues(dd.Queries)) == 0
}

func (dd QueryDIDDocRequestService) AreQueryValuesEmpty(c services.ResolverContext) bool {
	for _, v := range dd.Queries {
		// Queries is the map with list of string as value.
		// If there is only one value and it's empty string, then we need to return RepresentationNotSupported error
		if len(v) == 1 && v[0] == "" {
			return true
		}
	}
	return false
}

// function to register the all the query handlers
func (dd *QueryDIDDocRequestService) RegisterQueryHandlers(c services.ResolverContext) error {
	stopHandler := queries.StopHandler{}

	// Create Chain of responsibility
	// First we need to just ask for Did:

	// DidDoc handlers
	startHandler := diddocQueries.DidQueryAllVersionsHandler{}
	startHandler.IsDereferencing = dd.IsDereferencing
	lastHandler, err := dd.RegisterDidDocQueryHandlers(&startHandler, c)
	if err != nil {
		return err
	}

	if dd.AreResourceQueriesPlaced(c) {
		lastHandler, err := dd.RegisterResourceQueryHandlers(lastHandler, c)
		if err != nil {
			return err
		}
		err = lastHandler.SetNext(c, &stopHandler, dd.IsDereferencing)
		if err != nil {
			return err
		}
	} else {
		err = lastHandler.SetNext(c, &stopHandler, dd.IsDereferencing)
		if err != nil {
			return err
		}
	}

	dd.FirstHandler = &startHandler

	return nil
}

// function to register the DIDDoc query handlers
func (dd *QueryDIDDocRequestService) RegisterDidDocQueryHandlers(startHandler queries.BaseQueryHandlerI, c services.ResolverContext) (queries.BaseQueryHandlerI, error) {
	// - didQueryHandler
	// or
	// - versionIdHandler
	// After that we can find for service field if it's set.
	// VersionIdHandler -> VersionTimeHandler -> DidDocResolveHandler -> TransformKeysHandler -> DidDocMetadataHandler -> ServiceHandler -> RelativeRefHandler
	relativeRefHandler := diddocQueries.RelativeRefHandler{}
	serviceHandler := diddocQueries.ServiceHandler{}
	versionIdHandler := diddocQueries.VersionIdHandler{}
	versionTimeHandler := diddocQueries.VersionTimeHandler{}
	didDocResolveHandler := diddocQueries.DidDocResolveHandler{}
	transformKeysHandler := diddocQueries.TransformKeysHandler{}
	didDocMetadataHandler := diddocQueries.DidDocMetadataHandler{}

	err := startHandler.SetNext(c, &versionIdHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = versionIdHandler.SetNext(c, &versionTimeHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = versionTimeHandler.SetNext(c, &didDocResolveHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = didDocResolveHandler.SetNext(c, &transformKeysHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = transformKeysHandler.SetNext(c, &didDocMetadataHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = didDocMetadataHandler.SetNext(c, &serviceHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = serviceHandler.SetNext(c, &relativeRefHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	return &relativeRefHandler, nil
}

// function to register the Dereferencing or Resource query handlers
func (dd *QueryDIDDocRequestService) RegisterResourceQueryHandlers(startHandler queries.BaseQueryHandlerI, c services.ResolverContext) (queries.BaseQueryHandlerI, error) {
	// Resource handlers
	resourceQueryHandler := resourceQueries.ResourceQueryHandler{}
	resourceIdHandler := resourceQueries.ResourceIdHandler{}
	resourceMetadataHandler := resourceQueries.ResourceMetadataHandler{}
	resourceCollectionIdHandler := resourceQueries.ResourceCollectionIdHandler{}
	resourceNameHandler := resourceQueries.ResourceNameHandler{}
	resourceTypeHandler := resourceQueries.ResourceTypeHandler{}
	resourceVersionHandler := resourceQueries.ResourceVersionHandler{}
	resourceVersionTimeHandler := resourceQueries.ResourceVersionTimeHandler{}
	resourceValidationHandler := resourceQueries.ResourceValidationHandler{}
	resourceChecksumHandler := resourceQueries.ResourceChecksumHandler{}

	err := startHandler.SetNext(c, &resourceQueryHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	// It's a resource query to fetch the collection of resources
	err = resourceQueryHandler.SetNext(c, &resourceIdHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	// Resource handlers
	// Chain would be:
	// resourceQueryHandler -> resourceIdHandler -> resourceCollectionIdHandler ->
	// -> resourceNameHandler -> resourceTypeHandler -> resourceVersionHandler ->
	// -> resourceVersionTimeHandler -> resourceChecksumHandler -> resourceValidationHandler -> resourceMetadataHandler -> stopHandler
	err = resourceIdHandler.SetNext(c, &resourceCollectionIdHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceCollectionIdHandler.SetNext(c, &resourceNameHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceNameHandler.SetNext(c, &resourceTypeHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceTypeHandler.SetNext(c, &resourceVersionHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceVersionHandler.SetNext(c, &resourceChecksumHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceChecksumHandler.SetNext(c, &resourceVersionTimeHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceVersionTimeHandler.SetNext(c, &resourceValidationHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	err = resourceValidationHandler.SetNext(c, &resourceMetadataHandler, dd.IsDereferencing)
	if err != nil {
		return nil, err
	}

	return &resourceMetadataHandler, nil
}

func (dd *QueryDIDDocRequestService) Query(c services.ResolverContext) error {
	result, err := dd.FirstHandler.Handle(c, dd, nil)
	if err != nil {
		return err
	}
	if result == nil {
		return types.NewRepresentationNotSupportedError(dd.GetDid(), dd.GetContentType(), nil, dd.IsDereferencing)
	}
	return dd.SetResponse(result)
}

func (dd *QueryDIDDocRequestService) IsResourceData(result types.ResolutionResultI) bool {
	// If ContentStream is not DereferencedResourceData then it's not a resource data
	_result, ok := result.(*types.ResourceDereferencing)
	if !ok {
		return false
	}
	if _result.Metadata != nil {
		return false // Handle DereferenceResourceWithMetadata
	}
	_, ok = _result.ContentStream.(*types.DereferencedResourceData)
	return ok
}

func (dd QueryDIDDocRequestService) Respond(c services.ResolverContext) error {
	if dd.Result.IsRedirect() {
		return c.Redirect(http.StatusSeeOther, string(dd.Result.GetBytes()))
	}
	if dd.IsResourceData(dd.Result) {
		return dd.RespondWithResourceData(c)
	}
	return c.JSONPretty(http.StatusOK, dd.Result, "  ")
}
