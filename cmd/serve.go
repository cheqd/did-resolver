package cmd

import (
	"net/http"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	"github.com/cheqd/did-resolver/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
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

func serve() {
	log.Info().Msg("Loading configuration")
	config, err := utils.LoadConfig()
	if err != nil {
		panic(err)
	}

	log.Info().Msgf("Configuration: %s", config.MustMarshalJson())

	log.Info().Msgf("Setting log level: %s", config.LogLevel)
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)

	// Echo instance
	e := echo.New()
	e.HTTPErrorHandler = CustomHTTPErrorHandler

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Services
	ledgerService := services.NewLedgerService()

	mainnetEndpoint, err := utils.ParseGRPCEndpoint(config.MainnetEndpoint, "mainnet")
	if err != nil {
		panic(err)
	}
	testnetEndpoint, err := utils.ParseGRPCEndpoint(config.TestnetEndpoint, "testnet")
	if err != nil {
		panic(err)
	}
	networks := []types.Network{*mainnetEndpoint, *testnetEndpoint}

	for _, network := range networks {
		log.Info().Msgf("Registering network: %s.", network.Namespace)
		err := ledgerService.RegisterLedger(types.DID_METHOD, network)
		if err != nil {
			panic(err)
		}
	}

	requestService := services.NewRequestService(types.DID_METHOD, ledgerService)

	// Routes
	e.GET(types.RESOLVER_PATH+":did", requestService.ResolveDIDDoc)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource", requestService.DereferenceResourceData)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+":resource/metadata", requestService.DereferenceResourceMetadata)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+"all", requestService.DereferenceCollectionResources)
	e.GET(types.RESOLVER_PATH+":did"+types.RESOURCE_PATH+"", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "all")
	})

	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.ResolverListener))
}
