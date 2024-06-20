package timebucket

import (
	"cloud.google.com/go/bigquery"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/iterator"
	"strings"
	"time"
)

type TimeInterval string

const (
	HOUR TimeInterval = "HOUR"
	DAY  TimeInterval = "DAY"
)

type QueryInterval struct {
	Start time.Time
	End   time.Time
}

type BucketType struct {
	Interval   TimeInterval
	Multiplier int
}

type Bucket struct {
	StartTime time.Time
	Duration  time.Duration
}

type QueryServer struct {
	client        *bigquery.Client
	queryTable    *bigquery.TableMetadata
	queryLocation string
}

type UsageBucket struct {
	//DeviceUID, min as initial_reading, max as final_reading, (max - min) as kWh, readings, first_reading, last_reading, bucket
	DeviceUID    string
	Min          float64   `datastore:",omitempty"`
	Max          float64   `datastore:",omitempty"`
	Usage        float64   `datastore:",omitempty"`
	Readings     int       `datastore:",omitempty"`
	FirstReading time.Time `datastore:",omitempty"`
	LastReading  time.Time `datastore:",omitempty"`
	Bucket
}

func NewQueryServer(client *bigquery.Client, queryTable *bigquery.TableMetadata, queryLocation string) *QueryServer {
	return &QueryServer{
		client:        client,
		queryTable:    queryTable,
		queryLocation: queryLocation,
	}
}

func (dus QueryServer) readRow(r []bigquery.Value) (*UsageBucket, error) {
	if r == nil || len(r) < 7 {
		return nil, fmt.Errorf("invalid row")
	}
	ub := &UsageBucket{
		DeviceUID:    r[0].(string),
		Max:          r[1].(float64),
		Min:          r[2].(float64),
		Usage:        r[3].(float64),
		Readings:     int(r[4].(int64)),
		FirstReading: r[5].(time.Time),
		LastReading:  r[6].(time.Time),
		Bucket: Bucket{
			StartTime: r[7].(time.Time),
		},
	}

	return ub, nil
}

// RunQuery issue a query and show results.
func (dus QueryServer) RunQuery(bucket BucketType, interval *QueryInterval) ([]UsageBucket, error) {
	ctx := context.Background()

	//the FullID replacement is because of really terrible coding by google
	from := fmt.Sprintf("FROM `%s` ", strings.Replace(dus.queryTable.FullID, ":", ".", 1))
	if interval != nil {
		from = fmt.Sprintf(`%s WHERE Time > Timestamp("%s") AND time < Timestamp("%s")`, from, interval.Start.Format("2006-01-02 15:04:05"), interval.End.Format("2006-01-02 15:04:05"))
	}
	//very hard to get goland not to interpret this as sql hence the "SELECT " +
	query := "SELECT " +
		`DeviceUID, min as initial_reading, max as final_reading, (max - min) as kWh, readings, first_reading, last_reading, bucket from (` +
		fmt.Sprintf(`SELECT DeviceUID, Max(ReadingKWH) as max, Min(ReadingKWH) as min, Min(Time) as first_reading, Max(Time) as last_reading, count(*) as readings, 
							TIMESTAMP_BUCKET(Time, INTERVAL %d %s) AS bucket %s GROUP BY DeviceUID, bucket`, bucket.Multiplier, bucket.Interval, from) +
		") ORDER BY DeviceUID, bucket;"

	log.Debug().Str("query", query).Msg("about to run")

	q := dus.client.Query(query)
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = dus.queryLocation
	// Run the query and print results when the query job is completed.
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
	var usages []UsageBucket
	for {
		var row []bigquery.Value
		iErr := it.Next(&row)
		if iErr == iterator.Done {
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
		usages = append(usages, *ub)
		log.Info().Interface("row", row).Msg("new row")
	}
	return usages, nil
}
