package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/broker/vutility/helpers"
	"github.com/safecility/microservices/go/broker/vutility/server"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	jwtSecretName := fmt.Sprintf("projects/%s/secrets/jwt-key/versions/1", config.ProjectName)
	jwtSecret := setup.GetSecret(jwtSecretName)

	ctx := context.Background()
	psClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start service")
	}

	uplinksTopic := psClient.Topic(config.Topics.Uplinks)

	exists, err := uplinksTopic.Exists(ctx)
	if !exists || err != nil {
		log.Fatal().Err(err).Str("name", config.Topics.Uplinks).Msg("topic does not exist or error")
	}

	jwtParser := lib.NewJWTParser(jwtSecret)

	s := server.NewVutilityServer(&jwtParser, uplinksTopic)

	s.Start()
}
