package resources

// import "github.com/cheqd/did-resolver/types"

// func FilterResourceDereferencingByQueryValue(response types.ResolutionResultI, valuestring) (types.ResourceDereferencing, error) {
// 	// If response has type of ResourceDefereferencingResult,
// 	// then we need to check if the resourceCollectionId is the same as the one in the response
// 	resDeref, ok := response.(*types.ResourceDereferencing)
// 	if !ok {
// 		return nil, types.NewInternalError(service.GetDid(), service.GetContentType(), errors.New("ResourceCollectionIdHandler: response is not of type ResourceDereferencing"), d.IsDereferencing)
// 	}

// 	resourceCollection := resDeref.ContentStream.(*types.DereferencedResourceListStruct)
// 	resourceCollectionFiltered := resourceCollection.Resources.GetByResourceId(resourceCollectionId)
// 	if len(resourceCollectionFiltered) == 0 {
// 		return nil, types.NewNotFoundError(service.GetDid(), service.GetContentType(), nil, d.IsDereferencing)
// 	}

// 	resDeref.ContentStream = &types.DereferencedResourceListStruct{
// 		Resources: resourceCollectionFiltered,
// 	}

// 	return types.ResourceDereferencing{}, nil
// }
