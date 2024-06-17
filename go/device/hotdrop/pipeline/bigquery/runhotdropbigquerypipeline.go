package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/device/hotdrop/pipeline/bigquery/helpers"
	"github.com/safecility/microservices/go/device/hotdrop/pipeline/bigquery/server"
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

	sub := gpsClient.Subscription(config.Subscriptions.Hotdrop)
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failure on subscription exists")
	}
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Hotdrop).Msg("Subscription does not exist")
	}

	topic := gpsClient.Topic(config.Topics.Bigquery)
	exists, err = topic.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failure on topic exists")
	}
	if !exists {
		log.Fatal().Str("topic", config.Topics.Bigquery).Msg("Topic does not exist")
	}

	bigqueryServer := server.NewHotDropServer(sub, topic, config.StoreAll)

	bigqueryServer.Start()

}
