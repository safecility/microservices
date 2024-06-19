package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/messages"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/protobuffer"
	"net/http"
	"os"
)

type UsageServer struct {
	sub      *pubsub.Subscription
	pub      *pubsub.Topic
	encoding pubsub.SchemaEncoding
	storeAll bool
}

func NewHotDropServer(sub *pubsub.Subscription, pub *pubsub.Topic, storeAll bool) *UsageServer {
	return &UsageServer{sub: sub, pub: pub, storeAll: storeAll, encoding: pubsub.EncodingBinary}
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
			log.Debug().Str("eui", r.DeviceUID).Msg("skipping message as no device and storeAll == false")
			return
		}

		go func() {
			m := protobuffer.CreateProtobufMessage(r)
			r, crr := stream.PublishProtoToTopic(m, es.encoding, es.pub)
			if crr != nil {
				log.Err(crr).Msg("could not add milesight data")
			}
			log.Debug().Str("result", *r).Msg("published milesight data")
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
		port = "8085"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
