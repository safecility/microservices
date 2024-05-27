package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/process/hotdrop/messages"
	"github.com/safecility/microservices/go/process/hotdrop/store"
	"net/http"
	"os"
)

type HotDropServer struct {
	cache    store.DeviceStore
	store    *store.DatastoreHotdrop
	sub      *pubsub.Subscription
	pipeline *pubsub.Topic
}

func NewHotDropServer(cache store.DeviceStore, store *store.DatastoreHotdrop,
	sub *pubsub.Subscription, topic *pubsub.Topic) HotDropServer {

	return HotDropServer{sub: sub, cache: cache, store: store, pipeline: topic}
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

			log.Debug().Str("eui", r.DeviceUID).Msg("hotdrop message")

			d, irr := es.cache.GetDevice(r.DeviceUID)
			if irr != nil {
				log.Err(irr).Str("uid", r.DeviceUID).Msg("could not get device")
			}

			//if we don't have an admin device we still save hotdrop messages
			go func() {
				crr := es.store.AddHotdropMessage(&r)
				if crr != nil {
					log.Err(crr).Msg("could not add hotdrop data")
				}
			}()

			// if we aren't tracking the device move to next
			if d == nil {
				continue
			}

			if err != nil {
				log.Err(err).Msg("could not publish usage to topic")
			}
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
