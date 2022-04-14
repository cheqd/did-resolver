package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/cheqd/cheqd-did-resolver/services"
	"github.com/cheqd/cheqd-did-resolver/types"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"

	"strings"
)

type Config struct {
	Networks map[string]string `yaml:"networks"`
}

func main() {
	didResolutionPath := flag.String("path", "/1.0/identifier/:did", "URL path with DID resolution endpoint")
	didResolutionPort := flag.String("port", ":1313", "The endpoint port with DID resolution")
	flag.Parse()

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
	e.GET(*didResolutionPath, func(c echo.Context) error {
		did := c.Param("did")
		// decode the paramater
		// did, err := url.QueryUnescape(did1)
		//if err != nil {
		//	return c.JSON(http.StatusBadRequest, map[string]string{})
		//}
		accept := strings.Split(c.Request().Header.Get("accept"), ";")[0]
		if strings.Contains(accept, types.ResolutionJSONLDType) {
			accept = types.ResolutionDIDJSONLDType
		} else {
			accept = types.ResolutionDIDJSONType
		}
		resolutionOption := map[string]string{"Accept": accept}
		e.StdLogger.Println("get did")
		responseBody, err := requestService.ProcessDIDRequest(did, resolutionOption)
		status := http.StatusOK
		if err != "" {
			// todo: defined a correct status
			status = http.StatusBadRequest
		}
		return c.JSONBlob(status, []byte(responseBody))

		//
		//// track the resolution
		//atomic.AddUint64(&rt.resolves, 1)
		//c.Response().Header().Set(echo.HeaderContentType, didLDJson)
		//
		//return c.JSON(http.StatusOK, rr)
	})

	// Start server
	e.Logger.Fatal(e.Start(*didResolutionPort))
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
