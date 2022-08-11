package utils

import (
	"fmt"
	"strings"

	"github.com/cheqd/did-resolver/types"
	"github.com/spf13/viper"
)

func LoadConfig() (types.Config, error) {
	viper.SetConfigFile("config.yaml")
	err := viper.ReadInConfig()

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("CHEQD_DID_RESOLVER")
	viper.AutomaticEnv()

	if err != nil {
		return types.Config{}, fmt.Errorf("error reading config.yaml: %s", err)
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
