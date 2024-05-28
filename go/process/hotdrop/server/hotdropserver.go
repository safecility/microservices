package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/microservices/go/process/hotdrop/messages"
	"github.com/safecility/microservices/go/process/hotdrop/store"
	"net/http"
	"os"
)

type HotDropServer struct {
	cache        store.DeviceStore
	store        *store.DatastoreHotdrop
	sub          *pubsub.Subscription
	pipeline     *pubsub.Topic
	storeHotdrop bool
}

func NewHotDropServer(cache store.DeviceStore, store *store.DatastoreHotdrop,
	sub *pubsub.Subscription, topic *pubsub.Topic, storeHotdrop bool) HotDropServer {

	return HotDropServer{sub: sub, cache: cache, store: store, pipeline: topic, storeHotdrop: storeHotdrop}
}

func (es *HotDropServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *HotDropServer) receive() {

	err := es.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		newPowerData := &messages.VuSensorMessage{}

		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, newPowerData)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}
		log.Debug().Int("dataSize", len(newPowerData.Data)).Msg("vu hotdrop readings")

		hdData := newPowerData.GetHotDropReadings()
		for _, r := range hdData {
			log.Debug().Str("eui", r.DeviceEUI).Msg("hotdrop message")
			var pd *messages.PowerDevice

			if es.cache != nil {
				pd, err = es.cache.GetDevice(r.DeviceEUI)
				if err != nil {
					log.Warn().Err(err).Str("uid", r.DeviceEUI).Msg("could not get device")
				}
				if pd != nil {
					r.PowerDevice = pd
				}
			}

			//if we don't have an admin device we can still save hotdrop messages
			if es.storeHotdrop {
				go func() {
					crr := es.store.AddHotdropMessage(&r)
					if crr != nil {
						log.Err(crr).Msg("could not add hotdrop data")
					}
				}()
			}

			// if we aren't tracking the device move to next
			if pd == nil {
				continue
			}

			u := messages.GetPowerUsage(r, pd)
			topic, err := stream.PublishToTopic(u, es.pipeline)
			if err != nil {
				log.Err(err).Msg("could not publish usage to topic")
				continue
			}
			log.Debug().Str("topic", *topic).Msg("published usage to topic")
		}
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (es *HotDropServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8089"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
