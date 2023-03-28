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
	config := utils.GetConfig()
	// Setup logger
	utils.SetupLogger(config)
	// Services
	ledgerService := services.NewLedgerService()
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
	e.HTTPErrorHandler = utils.CustomHTTPErrorHandler

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

func main() {
	serve()
}
