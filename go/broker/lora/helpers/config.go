package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
)

type Deployment string

const (
	Local      Deployment = "local"
	Test       Deployment = "test"
	Staging    Deployment = "staging"
	Production Deployment = "prod"
)

type MqttConfig struct {
	AppID    string `json:"appID"`
	Username string `json:"username"`
	Address  string `json:"address"`
}

type Config struct {
	ProjectName string `json:"projectName"`
	Mqtt        MqttConfig
	Topics      struct {
		Joins            string `json:"joins"`
		Uplinks          string `json:"uplinks"`
		Downlinks        string `json:"downlinks"`
		DownlinkReceipts string `json:"downlinkReceipts"`
	} `json:"topics"`
	Subscriptions struct {
		Downlinks string `json:"downlinks"`
	} `json:"subscriptions"`
}

// GetConfig for ttn the Username has the form: username = fmt.Sprintf("%s@ttn", p.AppID)
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
