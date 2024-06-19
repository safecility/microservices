package main

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/helpers"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/setup/sections"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	tmd, err := sections.CheckOrCreateBigqueryTable(config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating bigquery table")
	}
	sections.SetupPubsub(config, tmd)

}
