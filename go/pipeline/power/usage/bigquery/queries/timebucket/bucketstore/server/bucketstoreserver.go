package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/store"
	"net/http"
	"os"
	"time"
)

type BucketStoreServer struct {
	queryServer *timebucket.QueryServer
	store       *store.DatastoreBuckets
}

func NewBucketStoreServer(queryServer *timebucket.QueryServer, store *store.DatastoreBuckets) *BucketStoreServer {
	return &BucketStoreServer{queryServer: queryServer, store: store}
}

func (es *BucketStoreServer) Start() {
	es.serverHttp()
}

// storeQuery TODO parse requests and call
func (es *BucketStoreServer) storeQuery() (int, error) {
	t := timebucket.BucketType{
		Interval:   "",
		Multiplier: 0,
	}
	i := &timebucket.QueryInterval{
		Start: time.Time{},
		End:   time.Time{},
	}
	r, err := es.queryServer.RunQuery(t, i)
	if err != nil {
		return 0, err
	}
	err = es.store.AddBuckets(r, &t)
	if err != nil {
		return 0, err
	}
	return len(r), nil
}

func (es *BucketStoreServer) serverHttp() {
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
