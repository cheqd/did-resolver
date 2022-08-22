package cmd

import (
	"net/http"
	"strings"

	cheqdUtils "github.com/cheqd/cheqd-node/x/cheqd/utils"
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

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Services
	ledgerService := services.NewLedgerService(config.Ledger.Timeout, config.Ledger.UseTls)

	networks := strings.Split(config.Ledger.Networks, ";")
	for _, network := range networks {
		args := strings.Split(network, "=")
		name, url := args[0], args[1]

		log.Info().Msgf("Registering network. Name: %s, url: %s.", name, url)
		err := ledgerService.RegisterLedger(config.Resolver.Method, name, url)
		if err != nil {
			panic(err)
		}
	}

	requestService := services.NewRequestService(config.Resolver.Method, ledgerService)

	// Routes
	e.GET(config.Api.ResolverPath + ":did", func(c echo.Context) error {
		didUrl := c.Param("did")
		log.Debug().Msgf("DID: %s", didUrl)

		accept := c.Request().Header.Get(echo.HeaderAccept)
		log.Trace().Msgf("Accept: %s", accept)

		c.QueryParams()

		var requestedContentType types.ContentType

		resolutionResponse := requestService.ProcessDIDRequest(didUrl, types.ResolutionOption{Accept: requestedContentType})

		c.Response().Header().Set(echo.HeaderContentType, resolutionResponse.GetContentType())
		return c.JSONPretty(resolutionResponse.GetStatus(), resolutionResponse, "  ")
	})


	e.GET(config.Api.ResolverPath + ":did/resources/:resource", requestService.DereferenceResourceData)
	e.GET(config.Api.ResolverPath + ":did/resource/:resource/metadata", requestService.DereferenceResourceMetadata)
	e.GET(config.Api.ResolverPath + ":did/resources/all", requestService.DereferenceCollectionResources)
	e.GET(config.Api.ResolverPath + ":did/resources/", func(c echo.Context) error {
		return c.Redirect(http.StatusMovedPermanently, "all")
	})

	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.Api.Listener))
}
