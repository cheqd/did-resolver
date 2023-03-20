package utils

import (
	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func SetupLogger(config types.Config) {
	log.Info().Msgf("Setting log level: %s", config.LogLevel)
	level, err := zerolog.ParseLevel(config.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
}

func GetConfig() types.Config {
	log.Info().Msg("Loading configuration")
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	log.Info().Msgf("Configuration: %s", config.MustMarshalJson())
	return config
}
