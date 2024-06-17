package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/usagestore/messages"
	"net/http"
	"os"
)

type UsageServer struct {
	sub      *pubsub.Subscription
	storeAll bool
}

func NewUsageWindowServer(sub *pubsub.Subscription, storeAll bool) UsageServer {
	return UsageServer{sub: sub, storeAll: storeAll}
}

func (es *UsageServer) Start() {
	go es.receive()
	es.serverHttp()
}

func (es *UsageServer) receive() {

	err := es.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		r := &messages.MeterReading{}

		log.Debug().Str("data", fmt.Sprintf("%s", message.Data)).Msg("raw data")
		err := json.Unmarshal(message.Data, r)
		message.Ack()
		if err != nil {
			log.Err(err).Msg("could not unmarshall data")
			return
		}

		if r.Device == nil && es.storeAll == false {
			log.Debug().Msg("skipping message as no device and storeAll == false")
			return
		}

		go func() {
			// do time windows for usage
		}()
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (es *UsageServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8084"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
