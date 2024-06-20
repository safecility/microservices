package main

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/helpers"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/server"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/store"
	"os"
)

func main() {

	ctx := context.Background()

	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	deployment = fmt.Sprintf("./bucketstore/%s", deployment)
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

	du := timebucket.NewQueryServer(client, tableMetadata, config.BigQuery.Region)
	//
	//qi := &timebucket.QueryInterval{
	//	Start: time.Now().Add(-time.Hour * 48),
	//	End:   time.Now(),
	//}
	//
	//tb := timebucket.BucketType{
	//	Interval:   timebucket.HOUR,
	//	Multiplier: 1,
	//}

	dsClient, err := datastore.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not get datastore client")
	}

	bStore := store.NewBucketDatastore(dsClient)

	bs := server.NewBucketStoreServer(du, bStore)

	bs.Start()

}
