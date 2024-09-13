package server

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/accumulator/helpers"
	"net/http"
	"os"
	"time"
)

type AccumulatorServer struct {
	queryServer *QueryGenerator
	client      *pubsub.Client
	config      helpers.PubsubConfig
}

func NewAccumulatorServer(queryServer *QueryGenerator, client *pubsub.Client, config helpers.PubsubConfig) *AccumulatorServer {
	return &AccumulatorServer{queryServer: queryServer, client: client, config: config}
}

func (as *AccumulatorServer) Start() {
	as.serverHttp()
}

func (as *AccumulatorServer) getAccumulatorTopic(accumulator string) (*pubsub.Topic, error) {
	ctx := context.Background()
	topicName := helpers.GetTopicName(accumulator, as.config)

	topic := as.client.Topic(topicName)
	e, err := topic.Exists(ctx)
	if err != nil {
		return nil, err
	}
	if !e {
		return nil, fmt.Errorf("topic doesn't exist")
	}
	return topic, nil
}

// runQuery TODO parse requests and call
func (as *AccumulatorServer) runQuery(accumulator string) error {
	topic, err := as.getAccumulatorTopic(accumulator)
	if err != nil {
		return err
	}
	t := gbigquery.BucketType{
		Interval:   "",
		Multiplier: 0,
	}
	i := &gbigquery.QueryInterval{
		Start: time.Time{},
		End:   time.Time{},
	}
	err = as.queryServer.RunPowerUsageQuery(accumulator, t, i, topic)

	return err
}

func (as *AccumulatorServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	http.HandleFunc("/hours", as.handleHours)
	http.HandleFunc("/hours/previous", as.handlePreviousHour)

	http.HandleFunc("/days", as.handleDays)
	http.HandleFunc("/days/previous", as.handlePreviousDay)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8086"
	}
	log.Debug().Msg(fmt.Sprintf("starting http server port %s", port))
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
