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
	e.GET(config.Api.ResolverPath, func(c echo.Context) error {
		didUrl := c.Param("did")
		log.Debug().Msgf("DID: %s", didUrl)

		accept := c.Request().Header.Get(echo.HeaderAccept)
		log.Trace().Msgf("Accept: %s", accept)

		var requestedContentType types.ContentType

		if strings.Contains(accept, "*/*") || strings.Contains(accept, string(types.DIDJSONLD)) {
			requestedContentType = types.DIDJSONLD
		} else if strings.Contains(accept, string(types.DIDJSON)) {
			requestedContentType = types.DIDJSON
		} else if strings.Contains(accept, string(types.JSONLD)) {
			requestedContentType = types.JSONLD
		} else {
			requestedContentType = types.JSON
		}
		log.Debug().Msgf("Requested content type: %s", requestedContentType)

		_, path, _, _, _ := cheqdUtils.TrySplitDIDUrl(didUrl)
		log.Debug().Msg(path)
		if utils.IsCollectionResourcesPathRedirect(path) {
			return c.Redirect(http.StatusMovedPermanently, "all")
		}
		resolutionResponse := requestService.ProcessDIDRequest(didUrl, types.ResolutionOption{Accept: requestedContentType})

		c.Response().Header().Set(echo.HeaderContentType, resolutionResponse.GetContentType())

		// if contentType != dereferencingOptions.Accept {
		// 	return didDereferencing.ContentStream, statusCode, contentType
		// }
		if utils.IsResourceDataPath(path) && resolutionResponse.GetStatus() == http.StatusOK {
			return c.Blob(resolutionResponse.GetStatus(), resolutionResponse.GetContentType(), resolutionResponse.GetBytes())
		}
		return c.JSONPretty(resolutionResponse.GetStatus(), resolutionResponse, "  ")
	})

	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.Api.Listener))
}
