package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usagestore/helpers"
	"github.com/safecility/microservices/go/pipeline/power/usagestore/server"
	"github.com/safecility/microservices/go/pipeline/power/usagestore/store"
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
	defer func(gpsClient *pubsub.Client) {
		err := gpsClient.Close()
		if err != nil {
			log.Err(err).Msg("Error closing pubsub client")
		}
	}(gpsClient)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return // this is here so golang doesn't complain about gpsClient being possibly nil
	}

	usageSubscription := gpsClient.Subscription(config.Subscriptions.Usage)
	exists, err := usageSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Usage).Msg("no subscription found")
	}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore client")
	}
	d, err := store.NewDatastoreUsage(dsClient, config.Store.Entity)

	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore for usage")
	}

	usageServer := server.NewUsageServer(d, usageSubscription, config.Store.StoreAll)
	usageServer.Start()
}
