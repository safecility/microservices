package server

import (
	"cloud.google.com/go/bigquery"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/accumulator/messages"
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

type QueryServer struct {
	client        *bigquery.Client
	queryTable    *bigquery.TableMetadata
	queryLocation string
}

func NewQueryServer(client *bigquery.Client, queryTable *bigquery.TableMetadata, queryLocation string) *QueryServer {
	return &QueryServer{
		client:        client,
		queryTable:    queryTable,
		queryLocation: queryLocation,
	}
}

func (dus QueryServer) readRow(a string, r []bigquery.Value) (*messages.UsageBucket, error) {
	if r == nil || len(r) < 3 {
		return nil, fmt.Errorf("invalid row")
	}
	ub := &messages.UsageBucket{
		SystemUID: r[0].(string),
		Usage:     r[1].(float64),
		Bucket: gbigquery.Bucket{
			StartTime: r[2].(time.Time),
		},
		Accumulator: a,
	}

	return ub, nil
}

// RunPowerUsageQuery pass the required BucketType and QueryInterval, the query uses TIMESTAMP_BUCKET and finds max and min
// values within the buckets
func (dus QueryServer) RunPowerUsageQuery(accumulator string, bucket gbigquery.BucketType, interval *gbigquery.QueryInterval) ([]messages.UsageBucket, error) {
	ctx := context.Background()

	//the FullID replacement is because of really terrible api work by google
	from := fmt.Sprintf("`%s` ", strings.Replace(dus.queryTable.FullID, ":", ".", 1))
	if interval != nil {
		from = fmt.Sprintf(`%s WHERE Time > Timestamp("%s") AND time < Timestamp("%s") AND %s IS NOT NULL`,
			from,
			interval.Start.UTC().Format("2006-01-02 15:04:05"), interval.End.UTC().Format("2006-01-02 15:04:05"),
			accumulator)
	} else {
		from = fmt.Sprintf(`%s WHERE %s IS NOT NULL`, from, accumulator)
	}

	//very hard to get goland not to interpret this as sql hence the "SELECT " +
	query := "SELECT " +
		fmt.Sprintf(`%s,  SUM(kWh) as kWh, bucket from ( 
SELECT DeviceUID, %s, (max - min) as kWh, bucket from ( SELECT `, accumulator, accumulator) +
		fmt.Sprintf(`%s, DeviceUID, Max(ReadingKWH) as max, Min(ReadingKWH) as min, TIMESTAMP_BUCKET(Time, INTERVAL %d %s) AS bucket 
	FROM %s GROUP BY %s, DeviceUID, bucket`, accumulator, bucket.Multiplier, bucket.Interval, from, accumulator) +
		fmt.Sprintf(") ) GROUP BY %s, bucket ORDER BY bucket", accumulator)

	log.Debug().Str("query", query).Msg("about to run")

	q := dus.client.Query(query)
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = dus.queryLocation
	// Run the query
	job, err := q.Run(ctx)
	if err != nil {
		return nil, err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return nil, err
	}
	if err := status.Err(); err != nil {
		return nil, err
	}
	it, err := job.Read(ctx)
	if err != nil || it == nil {
		log.Warn().Err(err).Msg("error or iterator is nil")
		return nil, err
	}
	var usages []messages.UsageBucket
	ty := bucket.String()
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
		ub, iErr := dus.readRow(accumulator, row)
		if iErr != nil || ub == nil {
			log.Warn().Err(iErr).Msg("row error")
			continue
		}
		ub.Type = ty
		usages = append(usages, *ub)
		log.Info().Interface("row", row).Msg("new row")
	}
	return usages, nil
}
