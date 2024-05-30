package main

import (
	"cloud.google.com/go/pubsub"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/mqtt/lib"
	"github.com/safecility/go/mqtt/server"
	"github.com/safecility/go/setup"
	"github.com/safecility/microservices/go/broker/lora/helpers"
	"net/http"
	"os"
)

func main() {
	deployment, isSet := os.LookupEnv(helpers.OSDeploymentKey)
	if !isSet {
		deployment = string(setup.Local)
	}
	config := helpers.GetConfig(deployment)

	ctx := context.Background()

	gpsClient, err := pubsub.NewClient(ctx, config.ProjectName)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
	}
	if gpsClient == nil {
		log.Fatal().Err(err).Msg("Failed to create pubsub client")
		return
	}

	uplinkTopic := gpsClient.Topic(config.Topics.Uplinks)
	exists, err := uplinkTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Uplinks).Msg("no uplink topic")
	}
	defer uplinkTopic.Stop()

	sub := gpsClient.Subscription(config.Subscriptions.Downlinks)
	exists, err = sub.Exists(ctx)
	if !exists {
		log.Fatal().Str("subscription", config.Subscriptions.Downlinks).Msg("no downlink subscription")
	}

	joinsTopic := gpsClient.Topic(config.Topics.Joins)
	exists, err = joinsTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.Joins).Msg("no join topic")
	}
	defer joinsTopic.Stop()

	downlinkReceiptsTopic := gpsClient.Topic(config.Topics.DownlinkReceipts)
	exists, err = downlinkReceiptsTopic.Exists(ctx)
	if !exists {
		log.Fatal().Str("topic", config.Topics.DownlinkReceipts).Msg("no downlinkReceipts topic")
	}
	defer downlinkReceiptsTopic.Stop()

	daliPubSub := server.GooglePubSub{
		Joins:            joinsTopic,
		Uplinks:          uplinkTopic,
		Downlinks:        sub,
		DownlinkReceipts: downlinkReceiptsTopic,
	}

	secretsClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create secrets client")
	}
	secrets := setup.GetNewSecrets(config.ProjectName, secretsClient)
	appKey, err := secrets.GetSecret(config.Secret)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get app key")
	}

	pc := server.MqttProxyConfig{
		AppID:           config.Mqtt.AppID,
		Username:        config.Mqtt.Username,
		MqttAddress:     config.Mqtt.Address,
		AppKey:          string(appKey),
		CanDownlink:     true,
		GooglePubSub:    daliPubSub,
		Transformer:     lib.TtnV3{AppID: config.Mqtt.AppID, UidTransformer: helpers.AppIdUidTransformer{AppID: config.Mqtt.AppID}},
		PayloadAdjuster: helpers.SimpleDaliPayloadAdjuster{},
	}

	p, err := server.NewPahoProxy(pc)
	if err != nil {
		log.Fatal().Err(err).Msg("could not create paho proxy")
	}

	err = p.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("could not run lora proxy")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8093"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
