package server

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"net/http"
	"time"
)

func (bss *BucketStoreServer) handlePreviousHour(w http.ResponseWriter, r *http.Request) {
	end := time.Now().Truncate(time.Hour)
	start := end.Add(-time.Hour)

	t := gbigquery.BucketType{
		Interval:   gbigquery.HOUR,
		Multiplier: 1,
	}
	i := &gbigquery.QueryInterval{
		Start: start,
		End:   end,
	}
	qr, err := bss.queryServer.RunQuery(t, i)
	if err != nil {
		log.Err(err).Msg("could not run query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = bss.store.AddBuckets(qr, &t)
	if err != nil {
		log.Err(err).Msg("could not store query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (bss *BucketStoreServer) previousDay(w http.ResponseWriter, r *http.Request) {}
