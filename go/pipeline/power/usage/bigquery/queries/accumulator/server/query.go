package server

import (
	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/pubsub"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/go/lib/stream"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/accumulator/messages"
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

const bigTableTime = "2006-01-02 15:04:05"

type QueryGenerator struct {
	client        *bigquery.Client
	queryTable    *bigquery.TableMetadata
	queryLocation string
}

func NewQueryGenerator(client *bigquery.Client, queryTable *bigquery.TableMetadata, queryLocation string) *QueryGenerator {
	return &QueryGenerator{
		client:        client,
		queryTable:    queryTable,
		queryLocation: queryLocation,
	}
}

// AccumulatorUID, Max, Min, Usage
func (dus QueryGenerator) readRow(r []bigquery.Value) (*messages.AccumulatorBucket, error) {
	if r == nil || len(r) < 5 {
		return nil, fmt.Errorf("invalid row")
	}
	ub := &messages.AccumulatorBucket{
		AccumulatorUID: r[0].(string),
		Max:            r[1].(float64),
		Min:            r[2].(float64),
		Usage:          r[3].(float64),
		BucketStart:    r[4].(time.Time),
	}
	if ub.Max == 0 || ub.Min > ub.Max {
		log.Warn().Str("row", fmt.Sprintf("%+v", r)).Str("usage", fmt.Sprintf("%+v", ub)).Msg("max is too small")
	}
	if ub.Min == 0 {
		log.Info().Msg("min reading is zero - check device is new")
	}

	return ub, nil
}

// RunPowerUsageQuery pass the required BucketType and QueryInterval, the query uses TIMESTAMP_BUCKET and finds max and min
// values within the buckets
func (dus QueryGenerator) RunPowerUsageQuery(accumulator string, bucket gbigquery.BucketType, interval *gbigquery.QueryInterval, topic *pubsub.Topic) error {
	ctx := context.Background()

	//the FullID replacement is because of really terrible coding by google
	from := fmt.Sprintf("`%s` WHERE SystemType='%s'", strings.Replace(dus.queryTable.FullID, ":", ".", 1), accumulator)
	if interval != nil {
		from = fmt.Sprintf(`%s AND Time > Timestamp("%s") AND time < Timestamp("%s")`, from, interval.Start.UTC().Format(bigTableTime), interval.End.UTC().Format(bigTableTime))
	}

	query := "SELECT " +
		fmt.Sprintf(`SystemUID, Sum(max) as Max, Sum(min) as Min, Sum(kWh) as kWh, bucket FROM (
  SELECT SystemUID, SystemElement, max, min, (max - min) as kWh, readings, bucket FROM (
    SELECT SystemUID, SystemType, SystemElement, Max(ReadingKWH) as max, Min(ReadingKWH) as min, Max(Time) as last_reading, Min(Time) as first_reading, count(*) as readings, TIMESTAMP_BUCKET(Time, INTERVAL %d %s) AS bucket FROM %s GROUP BY SystemUID, SystemType, SystemElement, bucket
    ) GROUP BY SystemUID, SystemElement, max, min, kWh, readings, bucket
) GROUP BY SystemUID, bucket`, bucket.Multiplier, bucket.Interval, from)

	q := dus.client.Query(query)
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = dus.queryLocation
	// Run the query and use results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err = status.Err(); err != nil {
		return err
	}
	it, err := job.Read(ctx)
	if err != nil || it == nil {
		log.Warn().Err(err).Msg("error or iterator is nil")
		return err
	}

	for {
		var row []bigquery.Value
		iErr := it.Next(&row)
		if errors.Is(iErr, iterator.Done) {
			break
		}
		if iErr != nil {
			log.Warn().Err(iErr).Msg("row error")
			continue
		}
		ub, iErr := dus.readRow(row)
		if iErr != nil || ub == nil {
			log.Warn().Err(iErr).Msg("row error")
			continue
		}
		ub.Accumulator = accumulator
		ub.BucketInterval = bucket.Interval

		_, err = stream.PublishToTopic(ub, topic)
		if err != nil {
			log.Error().Err(err).Msg("could not publish phase")
			return err
		}
	}
	return nil
}
