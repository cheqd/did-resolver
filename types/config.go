package types

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Resolver ResolverConfig
	Api      ApiConfig

	LogLevel string
}

type ResolverConfig struct {
	Method        string
	Networks      string
	LedgerTimeout string
}

type ApiConfig struct {
	Listener     string
	ResolverPath string
}

func(c *Config) MarshalYaml() (string, error) {
	bytes, err := yaml.Marshal(c)
	return string(bytes), err
}

func (c *Config) MustMarshalYaml() string {
	res, err := c.MarshalYaml()
	if err != nil {
		panic(err)
	}

	return res
}

func(c *Config) MarshalJson() (string, error) {
	bytes, err := json.Marshal(c)
	return string(bytes), err
}

func (c *Config) MustMarshalJson() string {
	res, err := c.MarshalJson()
	if err != nil {
		panic(err)
	}

	return res
}