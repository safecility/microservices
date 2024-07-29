package server

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/accumulator/store"
	"net/http"
	"os"
	"time"
)

type BucketStoreServer struct {
	queryServer *QueryServer
	store       *store.DatastoreBuckets
}

func NewBucketStoreServer(queryServer *QueryServer, store *store.DatastoreBuckets) *BucketStoreServer {
	return &BucketStoreServer{queryServer: queryServer, store: store}
}

func (bss *BucketStoreServer) Start() {
	bss.serverHttp()
}

// storeQuery TODO parse requests and call
func (bss *BucketStoreServer) storeQuery(accumulator string) (int, error) {
	t := gbigquery.BucketType{
		Interval:   "",
		Multiplier: 0,
	}
	i := &gbigquery.QueryInterval{
		Start: time.Time{},
		End:   time.Time{},
	}
	r, err := bss.queryServer.RunPowerUsageQuery(accumulator, t, i)
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

	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		_, err := fmt.Fprintf(w, "started")
		if err != nil {
			log.Err(err).Msg(fmt.Sprintf("could write to http.ResponseWriter"))
		}
	})

	router.GET("/healthcheck", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.WriteHeader(http.StatusOK)
	})

	router.GET("/hours/previous/:accumulator", bss.handlePreviousHour)

	router.GET("/days/previous/:accumulator", bss.handlePreviousDay)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}
	log.Debug().Msg(fmt.Sprintf("starting http accumulator server port %s", port))
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not start http")
	}
}
