package diddoc

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/services/queries"
	diddocQueries "github.com/cheqd/did-resolver/services/queries/diddoc"
	resourceQueries "github.com/cheqd/did-resolver/services/queries/resources"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
)

type QueryDIDDocRequestService struct {
	services.BaseRequestService
	FirstHandler queries.BaseQueryHandlerI
}

func (dd *QueryDIDDocRequestService) Setup(c services.ResolverContext) error {
	dd.IsDereferencing = true
	return nil
}

func (dd *QueryDIDDocRequestService) SpecificValidation(c services.ResolverContext) error {
	_, err := url.QueryUnescape(dd.Did)
	if err != nil {
		return types.NewInvalidDIDUrlError(dd.Did, dd.RequestedContentType, err, dd.IsDereferencing)
	}

	diff := types.AllSupportedQueries.DiffWithUrlValues(dd.Queries)
	if len(diff) > 0 {
		return types.NewRepresentationNotSupportedError("Queries from list: "+strings.Join(diff, ","), dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	versionId := dd.GetQueryParam(types.VersionId)
	versionTime := dd.GetQueryParam(types.VersionTime)
	service := dd.GetQueryParam(types.ServiceQ)
	relativeRef := dd.GetQueryParam(types.RelativeRef)

	// Validation of query parameters
	if versionId != "" && versionTime != "" {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	// relativeRef should be only with service parameter also
	if relativeRef != "" && service == "" {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	if versionTime != "" {
		_, err := utils.ParseFromStringTimeToGoTime(versionTime)
		if err != nil {
			return types.NewRepresentationNotSupportedError(versionTime, dd.RequestedContentType, err, dd.IsDereferencing)
		}
	}

	// Validate that versionId is UUID
	if versionId != "" && !utils.IsValidUUID(versionId) {
		return types.NewInvalidDIDUrlError(versionId, dd.RequestedContentType, nil, dd.IsDereferencing)
	}

	return nil
}

func (dd *QueryDIDDocRequestService) SpecificPrepare(c services.ResolverContext) error {
	queryRaw, flag := services.PrepareQueries(c)
	queries, err := url.ParseQuery(queryRaw)
	if err != nil {
		return err
	}
	if flag != nil {
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}
	dd.Queries = queries

	// Register query handlers
	return dd.RegisterQueryHandlers(c)
}

func (dd *QueryDIDDocRequestService) RegisterQueryHandlers(c services.ResolverContext) error {

	stopHandler := queries.StopHandler{}

	// Create Chain of responsibility
	// First we need to just ask for Did:

	// DidDoc handlers
	startHandler := diddocQueries.DidQueryHandler{}
	lastHandler, err := dd.RegisterDidDocQueryHanlders(&startHandler, c)
	if err != nil {
		return err
	}

	if len(types.ResourceSupportedQueries.IntersectWithUrlValues(dd.Queries)) > 0 {
		lastHandler, err := dd.RegisterResourceQueryHandlers(lastHandler, c)
		if err != nil {
			return err
		}
		err = lastHandler.SetNext(c, &stopHandler)
		if err != nil {
			return err
		}
	} else {
		err = lastHandler.SetNext(c, &stopHandler)
		if err != nil {
			return err
		}
	}

	dd.FirstHandler = &startHandler

	return nil
}

func (dd *QueryDIDDocRequestService) RegisterDidDocQueryHanlders(startHandler queries.BaseQueryHandlerI, c services.ResolverContext) (queries.BaseQueryHandlerI, error) {
	// - didQueryHandler
	// or
	// - versionIdHandler
	// After that we can find for service field if it's set.
	// didQueryHandler -> versionIdHandler -> versionTimeHandler -> serviceHandler -> stopHandler
	relativeRefHandler := diddocQueries.RelativeRefHandler{}
	serviceHandler := diddocQueries.ServiceHandler{}
	versionIdHandler := diddocQueries.VersionIdHandler{}
	versionTimeHandler := diddocQueries.VersionTimeHandler{}

	err := startHandler.SetNext(c, &versionIdHandler)
	if err != nil {
		return nil, err
	}

	err = versionIdHandler.SetNext(c, &versionTimeHandler)
	if err != nil {
		return nil, err
	}

	err = versionTimeHandler.SetNext(c, &serviceHandler)
	if err != nil {
		return nil, err
	}

	err = serviceHandler.SetNext(c, &relativeRefHandler)
	if err != nil {
		return nil, err
	}
	return &relativeRefHandler, nil

}

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

	err := startHandler.SetNext(c, &resourceQueryHandler)
	if err != nil {
		return nil, err
	}

	// It's a resource query to fetch the collection of resources
	err = resourceQueryHandler.SetNext(c, &resourceIdHandler)
	if err != nil {
		return nil, err
	}

	// Resource handlers
	// Chain would be:
	// resourceIdHandler -> resourceVersionTimeHandler -> resourceCollectionIdHandler ->
	// -> resourceNameHandler -> resourceTypeHandler -> resourceVersionHandler ->
	// -> resourceValidationHandler -> resourceMetadataHandler -> stopHandler
	err = resourceIdHandler.SetNext(c, &resourceVersionTimeHandler)
	if err != nil {
		return nil, err
	}

	err = resourceVersionTimeHandler.SetNext(c, &resourceCollectionIdHandler)
	if err != nil {
		return nil, err
	}

	err = resourceCollectionIdHandler.SetNext(c, &resourceNameHandler)
	if err != nil {
		return nil, err
	}

	err = resourceNameHandler.SetNext(c, &resourceTypeHandler)
	if err != nil {
		return nil, err
	}

	err = resourceTypeHandler.SetNext(c, &resourceVersionHandler)
	if err != nil {
		return nil, err
	}

	err = resourceVersionHandler.SetNext(c, &resourceValidationHandler)
	if err != nil {
		return nil, err
	}

	err = resourceValidationHandler.SetNext(c, &resourceMetadataHandler)
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
		return types.NewRepresentationNotSupportedError(dd.Did, dd.RequestedContentType, nil, dd.IsDereferencing)
	}
	return dd.SetResponse(result)
}

func (dd QueryDIDDocRequestService) Respond(c services.ResolverContext) error {
	if dd.Result.IsRedirect() {
		return c.Redirect(http.StatusSeeOther, string(dd.Result.GetBytes()))
	}
	return c.JSONPretty(http.StatusOK, dd.Result, "  ")
}
