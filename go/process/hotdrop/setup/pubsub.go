package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/process/hotdrop/helpers"
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

	uplinkTopic := gpsClient.Topic(config.Topics.Uplinks)
	exists, err := uplinkTopic.Exists(ctx)
	if !exists {
		uplinkTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Uplinks)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create topic")
		}
		log.Info().Str("topic", uplinkTopic.String()).Msg("created topic")
	}

	deviceStateTopic := gpsClient.Topic(config.Topics.Pipeline)
	exists, err = deviceStateTopic.Exists(ctx)
	if !exists {
		deviceStateTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Pipeline)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create topic")
		}
		log.Info().Str("topic", deviceStateTopic.String()).Msg("created topic")
	}

	uSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err = uSubscription.Exists(ctx)
	if !exists {
		uTopic := gpsClient.Topic(config.Topics.Uplinks)
		exists, err = uTopic.Exists(ctx)
		if !exists {
			uTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Uplinks)
			if err != nil {
				log.Fatal().Err(err).Str("topic", config.Topics.Uplinks).Msg("setup could not create topic")
			}
			log.Info().Str("topic", uTopic.String()).Msg("created topic")
		}

		r, err := time.ParseDuration("1hr")
		if err != nil {
			log.Fatal().Err(err).Msg("could not parse duration")
		}
		subConfig := stream.GetDefaultSubscriptionConfig(uTopic, r)
		uSubscription, err = gpsClient.CreateSubscription(ctx, config.Subscriptions.Uplinks, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
	}

	log.Info().Msg("setup complete")
}
