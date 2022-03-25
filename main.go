package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	//"net/url"
	"strings"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/identifier/:did", func(c echo.Context) error {
		did := c.Param("did")
		// decode the paramater
		// did, err := url.QueryUnescape(did1)
		//if err != nil {
		//	return c.JSON(http.StatusBadRequest, map[string]string{})
		//}
		accept := strings.Split(c.Request().Header.Get("accept"), ";")[0]
		return c.String(http.StatusOK, accept+did)
		//opt := resolver.ResolutionOption{Accept: accept}
		//rr := resolver.ResolveRepresentation(conn, did, opt)
		//
		//// add universal resolver specific data:
		//rr.ResolutionMetadata.DidProperties = map[string]string{
		//	"method":           "cosmos",
		//	"methodSpecificId": strings.TrimPrefix(rr.Document.Id, DidPrefix),
		//}
		//
		//// track the resolution
		//atomic.AddUint64(&rt.resolves, 1)
		//c.Response().Header().Set(echo.HeaderContentType, didLDJson)
		//
		//return c.JSON(http.StatusOK, rr)
	})

	// Start server
	e.Logger.Fatal(e.Start(":1313"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
