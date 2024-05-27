package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/process/hotdrop/helpers"
	"github.com/safecility/microservices/go/process/hotdrop/server"
	"github.com/safecility/microservices/go/process/hotdrop/store"
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

	pipelineTopic := gpsClient.Topic(config.Topics.Pipeline)
	exists, err := pipelineTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Pipeline).Msg("no uplink topic")
	}
	defer pipelineTopic.Stop()

	uplinksSubscription := gpsClient.Subscription(config.Subscriptions.Uplinks)
	exists, err = uplinksSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Uplinks).Msg("no uplinks subscription")
	}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not start service")
	}
	d, err := store.NewDatastoreHotdrop(dsClient)

	hotDropServer := server.NewHotDropServer(nil, d, uplinksSubscription, pipelineTopic)
	hotDropServer.Start()

}
