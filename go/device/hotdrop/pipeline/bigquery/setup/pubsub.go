package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/device/hotdrop/pipeline/bigquery/helpers"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()
	sClient, err := pubsub.NewSchemaClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create schema client")
	}
	defer func(sClient *pubsub.SchemaClient) {
		err := sClient.Close()
		if err != nil {
			log.Error().Err(err).Msg("could not close schema client")
		}
	}(sClient)

	schema, err := sClient.Schema(ctx, config.Schema.Name, pubsub.SchemaViewFull)
	if err != nil || schema == nil {
		schema, err = createProtoSchema(sClient, config.Schema.Name, config.Schema.FilePath)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create schema")
		}
	}

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	bigqueryTopic := gpsClient.Topic(config.Topics.Bigquery)
	exists, err := bigqueryTopic.Exists(ctx)
	if !exists {
		bigqueryTopic, err = createBigqueryTopic(gpsClient, config.Topics.Bigquery, schema)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create bigquery topic")
		}
	}

	bigquerySubscription := gpsClient.Subscription(config.Topics.Bigquery)
	exists, err = bigquerySubscription.Exists(ctx)
	if !exists {
		err = createBigQuerySubscription(gpsClient, config.Subscriptions.Hotdrop, config.Table, bigqueryTopic)
		if err != nil {
			log.Fatal().Err(err).Msg("could not create bigquery subscription")
		}
	}
}

// createProtoSchema creates a schema resource from a schema proto file.
func createProtoSchema(client *pubsub.SchemaClient, schemaID, protoFile string) (*pubsub.SchemaConfig, error) {
	protoSource, err := os.ReadFile(protoFile)
	if err != nil {
		return nil, fmt.Errorf("error reading from file: %s", protoFile)
	}

	config := pubsub.SchemaConfig{
		Type:       pubsub.SchemaProtocolBuffer,
		Definition: string(protoSource),
	}

	ctx := context.Background()
	s, err := client.CreateSchema(ctx, schemaID, config)
	if err != nil {
		return nil, fmt.Errorf("CreateSchema: %w", err)
	}
	log.Debug().Str("schema", s.Name).Msg("Schema created")
	return s, nil
}

func createBigqueryTopic(client *pubsub.Client, topicName string, schema *pubsub.SchemaConfig) (*pubsub.Topic, error) {
	ctx := context.Background()
	bigqueryTopic, err := client.CreateTopicWithConfig(ctx, topicName, &pubsub.TopicConfig{
		SchemaSettings: &pubsub.SchemaSettings{
			Schema:   schema.Name,
			Encoding: pubsub.EncodingBinary,
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("setup could not create topic")
	}
	log.Info().Str("topic", bigqueryTopic.String()).Msg("created topic")

	return bigqueryTopic, nil
}

// createBigQuerySubscription creates a Pub/Sub subscription that exports messages to BigQuery.
func createBigQuerySubscription(client *pubsub.Client, subscriptionName, table string, topic *pubsub.Topic) error {
	ctx := context.Background()

	sub, err := client.CreateSubscription(ctx, subscriptionName, pubsub.SubscriptionConfig{
		Topic: topic,
		BigQueryConfig: pubsub.BigQueryConfig{
			Table:         table,
			WriteMetadata: true,
		},
	})
	if err != nil {
		return fmt.Errorf("client.CreateSubscription: %w", err)
	}
	log.Debug().Str("subscription", sub.ID()).Msg("Created BigQuery subscription")

	return nil
}
