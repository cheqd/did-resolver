package resources

import (
	"github.com/cheqd/did-resolver/types"
	"github.com/labstack/echo/v4"
)

func SetRoutes(e *echo.Echo) {
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource", ResourceDataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource/metadata", ResourceMetadataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_METADATA, ResourceCollectionEchoHandler)
}
