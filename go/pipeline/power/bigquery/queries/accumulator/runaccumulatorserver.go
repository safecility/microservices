package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/bigquery/queries/accumulator/helpers"
	"github.com/safecility/microservices/go/pipeline/power/bigquery/queries/accumulator/server"
	"os"
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

	qs := server.NewQueryGenerator(client, tableMetadata, config.BigQuery.Region)

	psClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if psClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return
	}

	bs := server.NewAccumulatorServer(qs, psClient, config.Pubsub)

	bs.Start()

}
