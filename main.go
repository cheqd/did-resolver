package main

import (
	"fmt"
	"net/http"

	"github.com/cheqd/cheqd-did-resolver/services"
	"github.com/cheqd/cheqd-did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"strings"
)

func main() {
	viper.SetConfigFile("config.yaml")
	err1 := viper.ReadInConfig()

	viper.SetConfigFile(".env")
	err2 := viper.ReadInConfig()
	if err1 != nil && err2 != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %w \n Fatal error config file: %w\n", err1, err2))
	}

	viper.AutomaticEnv()

	didResolutionPath := viper.GetString("path")
	didResolutionListener := viper.GetString("listener")

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//setup
	networks := viper.GetStringMapString("networks")
	ledgerService := services.NewLedgerService()
	for network, url := range networks {
		e.StdLogger.Println(network)
		ledgerService.RegisterLedger(network, url)
	}
	requestService := services.NewRequestService(ledgerService)

	// Routes
	e.GET(didResolutionPath, func(c echo.Context) error {
		did := c.Param("did")
		accept := strings.Split(c.Request().Header.Get("accept"), ";")[0]
		if strings.Contains(accept, types.ResolutionJSONLDType) {
			accept = types.ResolutionDIDJSONLDType
		} else {
			accept = types.ResolutionDIDJSONType
		}
		e.StdLogger.Println("get did")
		responseBody, err := requestService.ProcessDIDRequest(did, types.ResolutionOption{Accept: accept})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Response().Header().Set(echo.HeaderContentType, accept)
		return c.JSONBlob(http.StatusOK, []byte(responseBody))
	})

	e.Logger.Fatal(e.Start(didResolutionListener))
}
