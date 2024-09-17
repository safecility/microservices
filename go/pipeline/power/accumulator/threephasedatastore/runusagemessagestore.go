package main

import (
	"cloud.google.com/go/datastore"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/helpers"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/server"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/store"
	"os"
	"time"
)

func main() {

	ctx := context.Background()

	var deployment string
	args := os.Args
	if len(args) == 2 {
		deployment = args[1]
	} else {
		var isSet bool
		deployment, isSet = os.LookupEnv(helpers.OSDeploymentKey)
		if !isSet {
			deployment = string(setup.Local)
		}
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

	usageSubscription := gpsClient.Subscription(config.Pubsub.Subscriptions.Usage)
	exists, err := usageSubscription.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Pubsub.Subscriptions.Usage).Msg("no subscription found")
	}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore client")
	}
	d := store.NewThreePhaseDatastore(dsClient, config.Store.Datastore.Entity)

	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore for usage")
	}

	id := config.Store.Firestore.Database

	fsClient, err := firestore.NewClientWithDatabase(ctx, config.ProjectName, id)
	f := store.NewDeviceFirestore(fsClient, int(time.Second))

	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore for usage")
	}

	usageServer := server.NewThreePhaseServer(usageSubscription, d, f)
	usageServer.Start()
}
