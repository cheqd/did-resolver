package types

import (
	"encoding/json"
	"time"
)

type RawConfig struct {
	MainnetEndpoint  string `mapstructure:"MAINNET_ENDPOINT"`
	TestnetEndpoint  string `mapstructure:"TESTNET_ENDPOINT"`
	ResolverListener string `mapstructure:"RESOLVER_LISTNER"`
	LogLevel         string `mapstructure:"LOG_LEVEL"`
}

type Config struct {
	Networks         []Network
	ResolverListener string
	LogLevel         string
}

type Network struct {
	Namespace string
	Endpoint  string
	UseTls    bool
	Timeout   time.Duration
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
