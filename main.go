package main

import (
	"net/http"
	"os"

	"github.com/cheqd/cheqd-did-resolver/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"

	//"net/url"
	"strings"
)

type Config struct {
	Networks map[string]string `yaml:"networks"`
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//setup
	e.StdLogger.Println("get config")
	config, err := getConfig("config.yml")
	e.StdLogger.Println(config)
	if err != nil {
		e.Logger.Fatal(err)
	}
	ledgerService := services.NewLedgerService()
	for network, url := range config.Networks {
		e.StdLogger.Println(network)
		ledgerService.RegisterLedger(network, url)
	}
	requestService := services.NewRequestService(ledgerService)

	// Routes
	e.GET("/identifier/:did", func(c echo.Context) error {
		did := c.Param("did")
		// decode the paramater
		// did, err := url.QueryUnescape(did1)
		//if err != nil {
		//	return c.JSON(http.StatusBadRequest, map[string]string{})
		//}
		accept := strings.Split(c.Request().Header.Get("accept"), ";")[0]
		resolutionOption := map[string]string{"Accept": accept}
		e.StdLogger.Println("get did")
		responseBody, err := requestService.ProcessDIDRequest(did, resolutionOption)
		status := http.StatusOK
		if err != "" {
			status = http.StatusBadRequest
		}
		return c.String(status, responseBody)
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

func getConfig(configFileName string) (Config, error) {
	f, err := os.Open(configFileName)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}
