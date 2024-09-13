package server

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"net/http"
	"time"
)

func (as *AccumulatorServer) handlePreviousHour(w http.ResponseWriter, _ *http.Request) {
	end := time.Now().Truncate(time.Hour)
	start := end.Add(-time.Hour)

	ty := gbigquery.BucketType{
		Interval:   gbigquery.HOUR,
		Multiplier: 1,
	}
	in := &gbigquery.QueryInterval{
		Start: start,
		End:   end,
	}
	topic, err := as.getAccumulatorTopic("phase")
	if err != nil {
		log.Err(err).Msg("could not get accumulator topic")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = as.queryServer.RunPowerUsageQuery("threePhase", ty, in, topic)
	if err != nil {
		log.Err(err).Msg("could not run query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (as *AccumulatorServer) handleHours(w http.ResponseWriter, _ *http.Request) {

	ty := gbigquery.BucketType{
		Interval:   gbigquery.HOUR,
		Multiplier: 1,
	}

	topic, err := as.getAccumulatorTopic("phase")
	if err != nil {
		log.Err(err).Msg("could not get accumulator topic")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = as.queryServer.RunPowerUsageQuery("threePhase", ty, nil, topic)
	if err != nil {
		log.Err(err).Msg("could not run query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (as *AccumulatorServer) handlePreviousDay(w http.ResponseWriter, _ *http.Request) {
	day := time.Hour * 24
	end := time.Now().Truncate(day)
	start := end.Add(-day)

	topic, err := as.getAccumulatorTopic("phase")
	if err != nil {
		log.Err(err).Msg("could not get accumulator topic")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ty := gbigquery.BucketType{
		Interval:   gbigquery.DAY,
		Multiplier: 1,
	}
	in := &gbigquery.QueryInterval{
		Start: start,
		End:   end,
	}

	err = as.queryServer.RunPowerUsageQuery("threePhase", ty, in, topic)
	if err != nil {
		log.Err(err).Msg("could not run query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (as *AccumulatorServer) handleDays(w http.ResponseWriter, r *http.Request) {

	topic, err := as.getAccumulatorTopic("threePhase")
	if err != nil {
		log.Err(err).Msg("could not get accumulator topic")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ty := gbigquery.BucketType{
		Interval:   gbigquery.DAY,
		Multiplier: 1,
	}

	err = as.queryServer.RunPowerUsageQuery("phase", ty, nil, topic)
	if err != nil {
		log.Err(err).Msg("could not run query")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (as *AccumulatorServer) previousDay(w http.ResponseWriter, r *http.Request) {}
