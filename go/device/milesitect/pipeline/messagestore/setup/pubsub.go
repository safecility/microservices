package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/device/milesitect/pipeline/messagestore/helpers"
	"os"
	"time"
)

func main() {

	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	milesiteSubscription := gpsClient.Subscription(config.Subscriptions.Milesite)
	exists, err := milesiteSubscription.Exists(ctx)
	if !exists {
		milesiteTopic := gpsClient.Topic(config.Topics.Milesite)
		exists, err = milesiteTopic.Exists(ctx)
		if !exists {
			milesiteTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Milesite)
			if err != nil {
				log.Fatal().Err(err).Str("topic", config.Topics.Milesite).Msg("setup could not create topic")
			}
			log.Info().Str("topic", milesiteTopic.String()).Msg("created topic")
		}

		r, err := time.ParseDuration("1h")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(milesiteTopic, r)
		milesiteSubscription, err = gpsClient.CreateSubscription(ctx, config.Subscriptions.Milesite, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
		log.Info().Str("topic", milesiteSubscription.String()).Msg("created subscription")
	}

	log.Info().Msg("setup complete")
}