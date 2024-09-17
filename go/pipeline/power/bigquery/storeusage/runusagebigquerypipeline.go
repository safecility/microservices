package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/bigquery/store/helpers"
	"github.com/safecility/microservices/go/pipeline/power/bigquery/store/server"
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

	sub := gpsClient.Subscription(config.Pubsub.Subscriptions.Usage)
	exists, err := sub.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Str("sub", config.Pubsub.Subscriptions.Usage).Msg("Failure on subscription exists")
	}
	if !exists {
		log.Fatal().Str("subscription", config.Pubsub.Subscriptions.Usage).Msg("Subscription does not exist")
	}

	topic := gpsClient.Topic(config.Pubsub.Topics.Bigquery)
	exists, err = topic.Exists(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failure on topic exists")
	}
	if !exists {
		log.Fatal().Str("topic", config.Pubsub.Topics.Bigquery).Msg("Topic does not exist")
	}

	bigqueryServer := server.NewUsageServer(sub, topic, config.StoreAll)

	bigqueryServer.Start()

}
