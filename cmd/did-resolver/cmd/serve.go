package cmd

import (
	"github.com/cheqd/did-resolver/services"
	didDocServices "github.com/cheqd/did-resolver/services/diddoc"
	resourceServices "github.com/cheqd/did-resolver/services/resource"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	// Import Echo Swagger middleware
	echoSwagger "github.com/swaggo/echo-swagger"

	// Import generated Swagger docs
	_ "github.com/cheqd/did-resolver/docs"
)

func getServeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Runs resolver as a web server",
		Run: func(cmd *cobra.Command, args []string) {
			serve()
		},
	}
}

func setupLogger(config types.Config) {
	log.Info().Msgf("Setting log level: %s", config.LogLevel)
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
}

func getConfig() types.Config {
	log.Info().Msg("Loading configuration")
	config, err := utils.LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Configuration: %s", config.MustMarshalJson())
	return config
}

func SetRoutes(e *echo.Echo) {
	// Routes
	// Did docs
	e.GET(types.SWAGGER_PATH, echoSwagger.WrapHandler)
	e.GET(types.RESOLVER_PATH+":did", didDocServices.DidDocEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSION_PATH+":version", didDocServices.DidDocVersionEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSION_PATH+":version/metadata", didDocServices.DidDocVersionMetadataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_VERSIONS_PATH, didDocServices.DidDocAllVersionMetadataEchoHandler)
	// Resources
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource", resourceServices.ResourceDataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource/metadata", resourceServices.ResourceMetadataEchoHandler)
	e.GET(types.RESOLVER_PATH+":did"+types.DID_METADATA, resourceServices.ResourceCollectionEchoHandler)
}

func serve() {
	// Get Config
	config := getConfig()
	// Setup logger
	setupLogger(config)
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
	e.HTTPErrorHandler = CustomHTTPErrorHandler

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
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	SetRoutes(e)

	e.Debug = true
	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.ResolverListener))
}
