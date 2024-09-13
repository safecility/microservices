package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

const (
	OSDeploymentKey = "DEPLOYMENT"
)

type PubsubConfig struct {
	Prefix string
	Suffix string
	Setup  []string
}

type Config struct {
	ProjectName string `json:"projectName"`
	BigQuery    struct {
		Dataset string `json:"dataset"`
		Table   string `json:"table"`
		Region  string `json:"region"`
	} `json:"bigQuery"`
	Pubsub PubsubConfig
}

func GetTopicName(accumulator string, config PubsubConfig) string {
	return fmt.Sprintf("%s-%s-%s", config.Prefix, accumulator, config.Suffix)
}

// GetConfig creates a config for the specified deployment
func GetConfig(deployment string) *Config {
	fileName := fmt.Sprintf("%s-config.json", deployment)

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not find config file")
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Err(err).Msg("could not close config during defer")
		}
	}(file)
	decoder := json.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal().Err(err).Str("filename", fileName).Msg("could not decode pubsub config")
	}
	return config
}
