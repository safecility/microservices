package main

import (
	"cloud.google.com/go/pubsub"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/device/eastron/process/helpers"
	"github.com/safecility/microservices/go/device/eastron/process/server"
	"github.com/safecility/microservices/go/device/eastron/process/store"
	"os"
)

func main() {

	ctx := context.Background()

	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return
	}

	uplinksSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err := uplinksSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Uplinks).Msg("no uplinks subscription")
	}

	eastronTopic := gpsClient.Topic(config.Topics.Eastron)
	exists, err = eastronTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Eastron).Msg("no hotdrop topic")
	}
	if err != nil {
		log.Fatal().Err(err).Str("topic", config.Topics.Eastron).Msg("could not get topic")
	}
	defer eastronTopic.Stop()

	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	defer func(secretsClient *secretmanager.Client) {
		err := secretsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close secrets client")
		}
	}(secretsClient)
	sqlSecret := setup.GetNewSecrets(config.ProjectName, secretsClient)
	password, err := sqlSecret.GetSecret(config.Sql.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get secret")
	}
	config.Sql.Config.Password = string(password)

	s, err := setup.NewSafecilitySql(config.Sql.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility sql")
	}
	c, err := store.NewDeviceSql(s)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup safecility device sql")
	}

	eagleServer := server.NewEastronServer(c, uplinksSubscription, eastronTopic, config.PipeAll)
	eagleServer.Start()

}
