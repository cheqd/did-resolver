package main

import (
	"github.com/cheqd/did-resolver/services"
	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	echoSwagger "github.com/swaggo/echo-swagger"

	// Import generated Swagger docs
	_ "github.com/cheqd/did-resolver/docs"
)

func serve() {
	// Get Config
	config := types.GetConfig()
	// Setup logger
	types.SetupLogger(config)
	
	// Initialize endpoint manager
	endpointManager := services.NewEndpointManager(config)
	
	// Services
	ledgerService := services.NewLedgerService(endpointManager)
	didService := services.NewDIDDocService(types.DID_METHOD, ledgerService)
	resourceService := services.NewResourceService(types.DID_METHOD, ledgerService)

	for _, network := range config.Networks {
		log.Info().Msgf("Registering network: %s.", network.Namespace)
		err := ledgerService.RegisterLedger(types.DID_METHOD, network)
		if err != nil {
			panic(err)
		}
	}

	// Echo instance
	e := echo.New()
	e.HTTPErrorHandler = services.CustomHTTPErrorHandler

	// Middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := services.ResolverContext{
				Context:         c,
				LedgerService:   ledgerService,
				DidDocService:   didService,
				ResourceService: resourceService,
			}
			return next(cc)
		}
	})

	// Client sends the Accept-Encoding header and
	// server should respond with the Content-Encoding header
	// Decompress only if gzip in headers
	e.Use(middleware.Decompress())

	// Compress only if gzip in headers
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		// If gzip not in Accept-Encoding header, do not compress
		Skipper: utils.GzipSkipper,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET(types.SWAGGER_PATH, echoSwagger.WrapHandler)

	didDocServices.SetRoutes(e)
	resourceServices.SetRoutes(e)

	e.Debug = true
	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.ResolverListener))
}

//	@title			DID Resolver for cheqd DID method
//	@version		v3.0
//	@description	Universal Resolver driver for cheqd DID method
//	@contact.name	Cheqd Foundation Limited
//	@contact.url	https://cheqd.io
//	@license.name	Apache 2.0
//	@license.url	https://github.com/cheqd/did-resolver/blob/main/LICENSE
//	@host			resolver.cheqd.net
//	@BasePath		/1.0/identifiers
//	@schemes		https http

func main() {
	err := types.PrintConfig()
	if err != nil {
		panic(err)
	}
	serve()
}
