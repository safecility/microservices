package server

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/store"
	"net/http"
	"os"
	"time"
)

type BucketStoreServer struct {
	queryServer *gbigquery.QueryServer
	store       *store.DatastoreBuckets
}

func NewBucketStoreServer(queryServer *gbigquery.QueryServer, store *store.DatastoreBuckets) *BucketStoreServer {
	return &BucketStoreServer{queryServer: queryServer, store: store}
}

func (bss *BucketStoreServer) Start() {
	bss.serverHttp()
}

// storeQuery TODO parse requests and call
func (bss *BucketStoreServer) storeQuery() (int, error) {
	t := gbigquery.BucketType{
		Interval:   "",
		Multiplier: 0,
	}
	i := &gbigquery.QueryInterval{
		Start: time.Time{},
		End:   time.Time{},
	}
	r, err := bss.queryServer.RunPowerUsageQuery(t, i)
	if err != nil {
		return 0, err
	}
	err = bss.store.AddBuckets(r, &t)
	if err != nil {
		return 0, err
	}
	return len(r), nil
}

func (bss *BucketStoreServer) serverHttp() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	http.HandleFunc("/hours/previous", bss.handlePreviousHour)

	http.HandleFunc("/days/previous", bss.handlePreviousDay)

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
