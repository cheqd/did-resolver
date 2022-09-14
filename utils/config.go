package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheqd/did-resolver/types"
	"github.com/spf13/viper"
)

func LoadConfig() (types.Config, error) {
	if _, err := os.Stat("config.env"); err == nil {
		viper.SetConfigFile("config.env")
		err := viper.ReadInConfig()
		if err != nil {
			return types.Config{}, fmt.Errorf("error reading config.env: %v", err)
		}
	}
	viper.SetDefault("MAINNET_ENDPOINT", "")
	viper.SetDefault("TESTNET_ENDPOINT", "")
	viper.SetDefault("LOG_LEVEL", "")
	viper.SetDefault("RESOLVER_LISTNER", "")
	viper.AutomaticEnv()

	rawConf := &types.RawConfig{}
	err := viper.Unmarshal(rawConf)
	if err != nil {
		return types.Config{}, fmt.Errorf("unable to decode into config struct, %v", err)
	}
	conf, err := NewConfig(*rawConf)
	if err != nil {
		return types.Config{}, fmt.Errorf("invalid config parameter, %v", err)
	}
	return conf, nil
}

func MustLoadConfig() types.Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}

	return config
}

func ParseGRPCEndpoint(configEndpoint string, networkName string) (*types.Network, error) {
	config := strings.Split(configEndpoint, ",")
	if len(config) != 3 {
		return nil, fmt.Errorf(fmt.Sprintf("Endpoint config for %s is invalid: %s", networkName, configEndpoint))
	}
	useTls, err := strconv.ParseBool(config[1])
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("useTls value %s for %s endpoint is invalid", configEndpoint, networkName))
	}
	timeout, err := time.ParseDuration(config[2])
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Timeout value %s for %s endpoint is invalid", configEndpoint, networkName))
	}

	return &types.Network{
		Namespace: networkName,
		Endpoint:  config[0],
		UseTls:    useTls,
		Timeout:   timeout,
	}, nil
}

func NewConfig(rawConfig types.RawConfig) (types.Config, error) {
	mainnetEndpoint, err := ParseGRPCEndpoint(rawConfig.MainnetEndpoint, "mainnet")
	if err != nil {
		return types.Config{}, err
	}
	testnetEndpoint, err := ParseGRPCEndpoint(rawConfig.TestnetEndpoint, "testnet")
	if err != nil {
		return types.Config{}, err
	}
	return types.Config{
		Networks: []types.Network{*mainnetEndpoint, *testnetEndpoint},
		ResolverListener: rawConfig.ResolverListener,
		LogLevel: rawConfig.LogLevel,
	}, nil
}
