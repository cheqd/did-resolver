package types

import (
	"encoding/json"
	"time"
)

// EndpointRole represents the role of an endpoint
type EndpointRole string

const (
	EndpointRolePrimary   EndpointRole = "primary"
	EndpointRoleFallback  EndpointRole = "fallback"
)

// Endpoint represents a gRPC endpoint with its configuration
type Endpoint struct {
	URL      string
	UseTls   bool
	Timeout  time.Duration
	Role     EndpointRole
}

// Network represents a blockchain network with endpoint configuration
type Network struct {
	Namespace string
	Endpoints []Endpoint
	UseTls    bool
	Timeout   time.Duration
}

type RawConfig struct {
	MainnetEndpoint           string `mapstructure:"MAINNET_ENDPOINT"`
	TestnetEndpoint           string `mapstructure:"TESTNET_ENDPOINT"`
	MainnetEndpointFallback   string `mapstructure:"MAINNET_ENDPOINT_FALLBACK"`
	TestnetEndpointFallback   string `mapstructure:"TESTNET_ENDPOINT_FALLBACK"`
	EnableFallbackEndpoints   bool   `mapstructure:"ENABLE_FALLBACK_ENDPOINTS"`
	ResolverListener          string `mapstructure:"RESOLVER_LISTENER"`
	LogLevel                  string `mapstructure:"LOG_LEVEL"`
}

type Config struct {
	Networks                []Network
	EnableFallbackEndpoints bool
	ResolverListener        string
	LogLevel               	string
}

func (c *Config) MarshalJson() (string, error) {
	bytes, err := json.MarshalIndent(c, "", "  ")
	return string(bytes), err
}

func (c *Config) MustMarshalJson() string {
	res, err := c.MarshalJson()
	if err != nil {
		panic(err)
	}

	return res
}
