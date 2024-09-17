package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/messages"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/store"
	"net/http"
	"os"
)

type ThreePhaseServer struct {
	datastore *store.ThreePhaseDatastore
	firestore *store.DeviceFirestore
	sub       *pubsub.Subscription
}

func NewThreePhaseServer(sub *pubsub.Subscription, d *store.ThreePhaseDatastore, f *store.DeviceFirestore) *ThreePhaseServer {
	return &ThreePhaseServer{sub: sub, datastore: d, firestore: f}
}

func (tps *ThreePhaseServer) Start() {
	go tps.receive()
	tps.serverHttp()
}

func (tps *ThreePhaseServer) receive() {

	log.Info().Str("sub", tps.sub.String()).Msg("starting message reception")

	err := tps.sub.Receive(context.Background(), func(ctx context.Context, message *pubsub.Message) {
		r := &messages.AccumulatorBucket{}

		err := json.Unmarshal(message.Data, r)
		message.Ack()
		if err != nil {
			log.Err(err).Str("data", fmt.Sprintf("%s", message.Data)).Msg("could not unmarshall data")
			return
		}
		log.Debug().Str("reading", fmt.Sprintf("%+v", r)).Msg("received message")

		outputDevice, err := tps.firestore.GetDevice(r.AccumulatorUID)
		if err != nil {
			log.Err(err).Str("data", fmt.Sprintf("%s", message.Data)).Msg("could not get device")
			return
		}

		mr := &messages.MeterReading{
			DeviceUID:   outputDevice.DeviceUID,
			DeviceName:  outputDevice.DeviceName,
			DeviceTag:   outputDevice.DeviceTag,
			DeviceType:  string(outputDevice.DeviceType),
			CompanyUID:  outputDevice.CompanyUID,
			LocationUID: outputDevice.LocationUID,
			ReadingKWH:  r.Usage,
			Time:        r.BucketStart,
			Interval:    r.BucketInterval,
		}

		go func() {
			crr := tps.datastore.AddMeterReading(mr)
			if crr != nil {
				log.Err(crr).Msg("could not add meter reading")
			}
		}()
	})
	if err != nil {
		log.Err(err).Msg("could not receive from sub")
		return
	}
}

func (tps *ThreePhaseServer) serverHttp() {
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
