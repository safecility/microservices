package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/dailyusage/helpers"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/dailyusage/server"
	"os"
	"time"
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

	du := server.NewDailyUsageServer(client, tableMetadata, config.BigQuery.Region)

	qi := &server.QueryInterval{
		Start: time.Now().Add(-time.Hour * 48),
		End:   time.Now(),
	}

	err = du.RunQuery(qi)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to run query")
		return
	}

}
