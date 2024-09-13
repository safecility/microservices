package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/accumulator/helpers"
	"os"
)

func main() {
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

	ctx := context.Background()
	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("could not setup pubsub")
	}

	for _, tn := range config.Pubsub.Setup {
		var irr error
		topicName := helpers.GetTopicName(tn, config.Pubsub)
		processorTopic := gpsClient.Topic(topicName)
		exists, irr := processorTopic.Exists(ctx)
		if irr != nil {
			log.Error().Err(irr).Msg("could not check if topic exists")
		} else if !exists {
			processorTopic, irr = gpsClient.CreateTopic(ctx, topicName)
			if irr != nil {
				log.Error().Err(irr).Msg("could not create topic")
			}
			log.Info().Str("topicName", topicName).Msg("topic created")
		} else {
			log.Info().Str("topicName", topicName).Msg("topic exists")
		}
	}

}
