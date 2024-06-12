package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/device/milesitect/pipeline/messagestore/messages"
	"github.com/safecility/microservices/go/device/milesitect/pipeline/messagestore/store"
	"net/http"
	"os"
)

type MilesiteServer struct {
	store    *store.DatastoreMilesite
	sub      *pubsub.Subscription
	storeAll bool
}

func NewEastronServer(store *store.DatastoreMilesite, sub *pubsub.Subscription, storeAll bool) MilesiteServer {
	return MilesiteServer{sub: sub, store: store, storeAll: storeAll}
}

func (es *MilesiteServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *MilesiteServer) receive() {

	err := es.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		r := &messages.EastronSdmReading{}

		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, r)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}

		if r.Device == nil && es.storeAll == false {
			log.Debug().Str("uid", r.UID).Msg("skipping message as no device and storeAll == false")
			return
		}

		go func() {
			crr := es.store.AddEastronMessage(r)
			if crr != nil {
				log.Err(crr).Msg("could not add hotdrop data")
			}
		}()
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
		port = "8081"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
