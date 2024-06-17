package sections

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/device/hotdrop/pipeline/bigquery/helpers"
)

func CheckOrCreateBigqueryTable(config *helpers.Config) (*bigquery.TableMetadata, error) {
	ctx := context.Background()

	client, err := bigquery.NewClient(ctx, config.ProjectName)
	if err != nil {
		return nil, fmt.Errorf("bigquery.NewClient: %v", err)
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
		log.Error().Err(err).Msg("Failed to get table metadata")
	}
	if tableMetadata == nil {
		err = tableRef.Create(ctx, getTableMetadata(config.BigQuery.Table))
		if err != nil {
			return nil, err
		}
		log.Info().Msg("Created bigquery table")
		tableMetadata, err = tableRef.Metadata(ctx)
	}

	return tableMetadata, err
}

func getTableMetadata(name string) *bigquery.TableMetadata {
	sampleSchema := bigquery.Schema{
		{Name: "DeviceEUI", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "InstantaneousCurrent", Type: bigquery.FloatFieldType},
		{Name: "MaximumCurrent", Type: bigquery.FloatFieldType},
		{Name: "SecondsAgoForMaximumCurrent", Type: bigquery.FloatFieldType},
		{Name: "AccumulatedCurrent", Type: bigquery.FloatFieldType},
		{Name: "MinimumCurrent", Type: bigquery.FloatFieldType},
		{Name: "SecondsAgoForMinimumCurrent", Type: bigquery.FloatFieldType},
		{Name: "SupplyVoltage", Type: bigquery.FloatFieldType},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: sampleSchema,
	}
}
