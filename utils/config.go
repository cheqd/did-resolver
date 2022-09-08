package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cheqd/did-resolver/types"
	"github.com/spf13/viper"
)

func LoadConfig() (types.Config, error) {
	viper.SetConfigFile("config.env")
	

	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("error reading config.env: %s", err)
	}

	conf := &types.Config{}
	err = viper.Unmarshal(conf)
	if err != nil {
		return types.Config{}, fmt.Errorf("unable to decode into config struct, %v", err)
	}

	return *conf, nil
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
