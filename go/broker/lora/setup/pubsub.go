package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/broker/lora/dali/helpers"
	"os"
	"time"
)

func main() {

	deployment, isSet := os.LookupEnv("Deployment")
	if !isSet {
		deployment = string(helpers.Local)
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
	}
	downlinksSub := gpsClient.Subscription(config.Subscriptions.Downlinks)
	exists, err = downlinksSub.Exists(ctx)
	if !exists {
		downlinksTopic := gpsClient.Topic(config.Topics.Downlinks)
		exists, err = downlinksTopic.Exists(ctx)
		if !exists {
			downlinksTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Downlinks)
			if err != nil {
				log.Fatal().Err(err).Msg("setup could not create topic")
			}
		}

		subConfig := getSubscriptionConfig(downlinksTopic)
		downlinksSub, err = gpsClient.CreateSubscription(ctx, config.Subscriptions.Downlinks, subConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create subscription")
		}
	}

	joinsTopic := gpsClient.Topic(config.Topics.Joins)
	exists, err = joinsTopic.Exists(ctx)
	if !exists {
		joinsTopic, err = gpsClient.CreateTopic(ctx, config.Topics.Joins)
		if err != nil {
			log.Fatal().Err(err).Msg("setup could not create topic")
		}
	}

	downlinkReceiptsTopic := gpsClient.Topic(config.Topics.DownlinkReceipts)
	exists, err = downlinkReceiptsTopic.Exists(ctx)
	if !exists {
		downlinkReceiptsTopic, err = gpsClient.CreateTopic(ctx, config.Topics.DownlinkReceipts)
	}

	log.Info().Msg("finished setup")

}

func getSubscriptionConfig(topic *pubsub.Topic) pubsub.SubscriptionConfig {
	retentionDuration, err := time.ParseDuration("24h")
	if err != nil {
		log.Err(err).Msg("could not create duration")
		retentionDuration = 0
	}

	return pubsub.SubscriptionConfig{
		Topic:                         topic,
		PushConfig:                    pubsub.PushConfig{},
		BigQueryConfig:                pubsub.BigQueryConfig{},
		CloudStorageConfig:            pubsub.CloudStorageConfig{},
		AckDeadline:                   0,
		RetainAckedMessages:           false,
		RetentionDuration:             retentionDuration,
		ExpirationPolicy:              nil,
		Labels:                        nil,
		EnableMessageOrdering:         false,
		DeadLetterPolicy:              nil,
		Filter:                        "",
		RetryPolicy:                   nil,
		Detached:                      false,
		TopicMessageRetentionDuration: 0,
		EnableExactlyOnceDelivery:     false,
		State:                         0,
	}
}
