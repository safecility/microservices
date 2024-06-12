package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/microservices/go/device/milesitect/process/messages"
	"github.com/safecility/microservices/go/device/milesitect/process/store"
	"net/http"
	"os"
)

type MilesiteServer struct {
	cache         store.DeviceStore
	sub           *pubsub.Subscription
	milesiteTopic *pubsub.Topic
	pipeAll       bool
}

func NewMilesiteServer(cache store.DeviceStore, sub *pubsub.Subscription, eagle *pubsub.Topic, pipeAll bool) MilesiteServer {
	return MilesiteServer{sub: sub, cache: cache, milesiteTopic: eagle, pipeAll: pipeAll}
}

func (es *MilesiteServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *MilesiteServer) receive() {

	err := es.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		sm := &stream.SimpleMessage{}
		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, sm)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}

		mr, err := messages.ReadMilesiteCT(sm.Payload)
		if err != nil {
			log.Err(err).Msg("could not read milesite CT")
			return
		}

		log.Debug().Str("eui", sm.DeviceUID).Msg("eastron message")
		var pd *messages.PowerDevice
		if es.cache != nil {
			pd, err = es.cache.GetDevice(sm.DeviceUID)
			if err != nil {
				log.Warn().Err(err).Str("uid", sm.DeviceUID).Msg("could not get device")
			}
			if pd == nil {
				log.Debug().Str("uid", sm.DeviceUID).Msg("device not found")
			}
			mr.PowerDevice = pd
		}
		if mr.Device == nil && !es.pipeAll {
			log.Debug().Str("device", sm.DeviceUID).Msg("no device in cache and pipeAll == false")
			return
		}
		mr.Time = sm.Time

		topic, err := stream.PublishToTopic(mr, es.milesiteTopic)
		if err != nil {
			log.Err(err).Msg("could not publish usage to topic")
			return
		}
		log.Debug().Str("topic", *topic).Msg("published eastron sdm to topic")
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (es *MilesiteServer) serverHttp() {
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
