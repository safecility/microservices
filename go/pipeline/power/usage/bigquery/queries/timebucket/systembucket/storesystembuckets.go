package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/systembucket/helpers"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/systembucket/server"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/systembucket/store"
	"os"
)

func main() {

	ctx := context.Background()

	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	client, err := bigquery.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create client")
	}
	defer func(client *bigquery.Client) {
		err := client.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close bigquery.Client")
		}
	}(client)

	tableRef := client.Dataset(config.BigQuery.Dataset).Table(config.BigQuery.Table)
	tableMetadata, err := tableRef.Metadata(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get table metadata")
	}
	if tableMetadata == nil {
		log.Fatal().Err(err).Msg("Failed to get table metadata")
	}

	du := server.NewQueryServer(client, tableMetadata, config.BigQuery.Region)

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore client")
	}

	defer func(dsClient *datastore.Client) {
		err = dsClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("Failed to close datastore client")
		}
	}(dsClient)

	bStore := store.NewBucketDatastore(dsClient)

	bs := server.NewBucketStoreServer(du, bStore)

	bs.Start()

}
