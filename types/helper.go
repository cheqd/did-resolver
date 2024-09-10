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

func ParseGRPCEndpoint(configEndpoint string, networkName string) (*Network, error) {
	config := strings.Split(configEndpoint, ",")
	if len(config) != 3 {
		return nil, fmt.Errorf("endpoint config for %s is invalid: %s", networkName, configEndpoint)
	}
	useTls, err := strconv.ParseBool(config[1])
	if err != nil {
		return nil, fmt.Errorf("useTls value %s for %s endpoint is invalid", configEndpoint, networkName)
	}
	timeout, err := time.ParseDuration(config[2])
	if err != nil {
		return nil, fmt.Errorf("timeout value %s for %s endpoint is invalid", configEndpoint, networkName)
	}

	return &Network{
		Namespace: networkName,
		Endpoint:  config[0],
		UseTls:    useTls,
		Timeout:   timeout,
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
	mainnetEndpoint, err := ParseGRPCEndpoint(rawConfig.MainnetEndpoint, "mainnet")
	if err != nil {
		return Config{}, err
	}
	testnetEndpoint, err := ParseGRPCEndpoint(rawConfig.TestnetEndpoint, "testnet")
	if err != nil {
		return Config{}, err
	}
	return Config{
		Networks:         []Network{*mainnetEndpoint, *testnetEndpoint},
		ResolverListener: rawConfig.ResolverListener,
		LogLevel:         rawConfig.LogLevel,
	}, nil
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
