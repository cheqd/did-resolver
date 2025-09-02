package types

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func AddElemToSet(set []string, newElement string) []string {
	if set == nil {
		set = []string{}
	}
	for _, c := range set {
		if c == newElement {
			return set
		}
	}
	return append(set, newElement)
}

func SetupLogger(config Config) {
	log.Info().Msgf("Setting log level: %s", config.LogLevel)
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
}

func ParseGRPCEndpoint(configEndpoint string) (*Endpoint, error) {
	config := strings.Split(configEndpoint, ",")
	if len(config) != 3 {
		return nil, fmt.Errorf("endpoint config is invalid: %s", configEndpoint)
	}
	useTls, err := strconv.ParseBool(config[1])
	if err != nil {
		return nil, fmt.Errorf("useTls value %s is invalid", configEndpoint)
	}
	timeout, err := time.ParseDuration(config[2])
	if err != nil {
		return nil, fmt.Errorf("timeout value %s is invalid", configEndpoint)
	}

	return &Endpoint{
		URL:     config[0],
		UseTls:  useTls,
		Timeout: timeout,
	}, nil
}

// Config functions

func LoadConfig() (Config, error) {
	if _, err := os.Stat("config.env"); err == nil {
		viper.SetConfigFile("config.env")
		err := viper.ReadInConfig()
		if err != nil {
			return Config{}, fmt.Errorf("error reading config.env: %v", err)
		}
	}
	viper.SetDefault("MAINNET_ENDPOINT", "")
	viper.SetDefault("TESTNET_ENDPOINT", "")
	viper.SetDefault("MAINNET_ENDPOINT_FALLBACK", "")
	viper.SetDefault("TESTNET_ENDPOINT_FALLBACK", "")
	viper.SetDefault("ENABLE_FALLBACK_ENDPOINTS", false)
	viper.SetDefault("LOG_LEVEL", "")
	viper.SetDefault("RESOLVER_LISTENER", "")
	viper.AutomaticEnv()

	rawConf := &RawConfig{}
	err := viper.Unmarshal(rawConf)
	if err != nil {
		return Config{}, fmt.Errorf("unable to decode into config struct, %v", err)
	}
	conf, err := NewConfig(*rawConf)
	if err != nil {
		return Config{}, fmt.Errorf("invalid config parameter, %v", err)
	}
	return conf, nil
}

func MustLoadConfig() Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func NewConfig(rawConfig RawConfig) (Config, error) {
	// Parse primary endpoints
	mainnetPrimary, err := ParseGRPCEndpoint(rawConfig.MainnetEndpoint)
	if err != nil {
		return Config{}, err
	}
	testnetPrimary, err := ParseGRPCEndpoint(rawConfig.TestnetEndpoint)
	if err != nil {
		return Config{}, err
	}

	// Create networks with primary endpoints
	networks := []Network{
		{
			Namespace: "mainnet",
			Endpoints: []Endpoint{
				{
					URL:     mainnetPrimary.URL,
					UseTls:  mainnetPrimary.UseTls,
					Timeout: mainnetPrimary.Timeout,
					Role:    EndpointRolePrimary,
				},
			},
			UseTls:   mainnetPrimary.UseTls,
			Timeout:  mainnetPrimary.Timeout,
		},
		{
			Namespace: "testnet",
			Endpoints: []Endpoint{
				{
					URL:     testnetPrimary.URL,
					UseTls:  testnetPrimary.UseTls,
					Timeout: testnetPrimary.Timeout,
					Role:    EndpointRolePrimary,
				},
			},
			UseTls:   testnetPrimary.UseTls,
			Timeout:  testnetPrimary.Timeout,
		},
	}

	// Handle fallback endpoints if enabled
	if rawConfig.EnableFallbackEndpoints {
		// When fallbacks are enabled, ALL namespaces must have fallback endpoints
		if rawConfig.MainnetEndpointFallback == "" {
			return Config{}, fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but MAINNET_ENDPOINT_FALLBACK is not configured")
		}
		if rawConfig.TestnetEndpointFallback == "" {
			return Config{}, fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but TESTNET_ENDPOINT_FALLBACK is not configured")
		}
		
		// Parse fallback endpoints
		mainnetFallback, err := ParseGRPCEndpoint(rawConfig.MainnetEndpointFallback)
		if err != nil {
			return Config{}, fmt.Errorf("invalid mainnet fallback endpoint: %v", err)
		}
		testnetFallback, err := ParseGRPCEndpoint(rawConfig.TestnetEndpointFallback)
		if err != nil {
			return Config{}, fmt.Errorf("invalid testnet fallback endpoint: %v", err)
		}
		
		// Add fallback endpoints to existing networks by namespace
		for i, network := range networks {
			if network.Namespace == "mainnet" {
				networks[i].Endpoints = append(networks[i].Endpoints, Endpoint{
					URL:     mainnetFallback.URL,
					UseTls:  mainnetFallback.UseTls,
					Timeout: mainnetFallback.Timeout,
					Role:    EndpointRoleFallback,
				})
			} else if network.Namespace == "testnet" {
				networks[i].Endpoints = append(networks[i].Endpoints, Endpoint{
					URL:     testnetFallback.URL,
					UseTls:  testnetFallback.UseTls,
					Timeout: testnetFallback.Timeout,
					Role:    EndpointRoleFallback,
				})
			}
		}
		
		// Validate that each namespace has at least 2 endpoints (primary + fallback)
		if err := validateFallbackEndpoints(networks); err != nil {
			return Config{}, err
		}
	}

	return Config{
		Networks:                networks,
		EnableFallbackEndpoints: rawConfig.EnableFallbackEndpoints,
		ResolverListener:        rawConfig.ResolverListener,
		LogLevel:                rawConfig.LogLevel,
	}, nil
}

// validateFallbackEndpoints ensures that when fallbacks are enabled, each namespace has at least 2 endpoints
func validateFallbackEndpoints(networks []Network) error {
	if len(networks) == 0 {
		return fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but no fallback endpoints configured")
	}
	
	for _, network := range networks {
		if len(network.Endpoints) < 2 {
			return fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but namespace %s only has %d endpoint(s) (need at least 2: primary + fallback)", network.Namespace, len(network.Endpoints))
		}
		
		// Ensure both primary and fallback endpoints exist
		primaryFound := false
		fallbackFound := false
		for _, endpoint := range network.Endpoints {
			if endpoint.Role == EndpointRolePrimary {
				primaryFound = true
			}
			if endpoint.Role == EndpointRoleFallback {
				fallbackFound = true
			}
		}
		if !primaryFound {
			return fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but namespace %s missing primary endpoint", network.Namespace)
		}
		if !fallbackFound {
			return fmt.Errorf("ENABLE_FALLBACK_ENDPOINTS=true but namespace %s missing fallback endpoint", network.Namespace)
		}
	}
	
	return nil
}

func PrintConfig() error {
	config := MustLoadConfig()
	configJson := config.MustMarshalJson()

	println(configJson)

	return nil
}

func GetConfig() Config {
	log.Info().Msg("Loading configuration")
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Configuration: %s", config.MustMarshalJson())
	return config
}
