package cmd

import (
	"net/http"
	"strings"

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
		err := ledgerService.RegisterLedger(name, url)
		if err != nil {
			panic(err)
		}
	}

	requestService := services.NewRequestService(config.Resolver.Method, ledgerService)

	// Routes
	e.GET(config.Api.ResolverPath, func(c echo.Context) error {
		didUrl := c.Param("did")
		log.Debug().Msgf("DID: %s", didUrl)

		acceptSettings := strings.Split(c.Request().Header.Get(echo.HeaderAccept), ";")
		accept := acceptSettings[0]
		if len(acceptSettings) > 1 && strings.Contains(acceptSettings[1], "profile=") {
			accept += ";" + acceptSettings[1]
		}
		log.Trace().Msgf("Accept: %s", accept)

		requestedContentType := types.ContentType(accept)
		if accept == "*/*" || accept == string(types.JSONLD) {
			requestedContentType = types.DIDJSONLD
		} else if strings.Contains(accept, string(types.HTML)) {
			requestedContentType = types.HTML
		}
		log.Debug().Msgf("Requested content type: %s", requestedContentType)

		responseBody, err := requestService.ProcessDIDRequest(didUrl, types.ResolutionOption{Accept: requestedContentType})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		log.Debug().Msgf("Response body: %s", responseBody)

		c.Response().Header().Set(echo.HeaderContentType, string(requestedContentType))
		return c.JSONBlob(http.StatusOK, []byte(responseBody))
	})

	log.Info().Msg("Starting listener")
	log.Fatal().Err(e.Start(config.Api.Listener))
}
