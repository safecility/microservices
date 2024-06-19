package server

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type DailyUsageServer struct {
	client        *bigquery.Client
	queryTable    *bigquery.TableMetadata
	queryLocation string
}

type QueryInterval struct {
	Start time.Time
	End   time.Time
}

func NewDailyUsageServer(client *bigquery.Client, queryTable *bigquery.TableMetadata, queryLocation string) *DailyUsageServer {
	return &DailyUsageServer{
		client:        client,
		queryTable:    queryTable,
		queryLocation: queryLocation,
	}
}

// RunQuery issue a query and show results.
func (dus DailyUsageServer) RunQuery(interval *QueryInterval) error {
	ctx := context.Background()

	//the FullID replacement is because of really terrible coding by google
	from := fmt.Sprintf("FROM `%s` ", strings.Replace(dus.queryTable.FullID, ":", ".", 1))
	if interval != nil {
		from = fmt.Sprintf(`%s WHERE Time > Timestamp("%s") AND time < Timestamp("%s")`, from, interval.Start.Format("2006-01-02 15:04:05"), interval.End.Format("2006-01-02 15:04:05"))
	}
	//very hard to get goland not to interpret this as sql hench the +
	query := "SELECT " +
		`DeviceUID, min as initial_reading, max as final_reading, (max - min) as kWh, readings, first_reading, last_reading, bucket from (` +
		fmt.Sprintf(`SELECT DeviceUID, Max(ReadingKWH) as max, Min(ReadingKWH) as min, Min(Time) as first_reading, Max(Time) as last_reading, count(*) as readings, 
							TIMESTAMP_BUCKET(Time, INTERVAL 1 DAY) AS bucket %s GROUP BY DeviceUID, bucket`, from) +
		") ORDER BY DeviceUID, bucket;"

	log.Error().Str("query", query).Msg("query is")

	q := dus.client.Query(query)
	// Location must match that of the dataset(s) referenced in the query.
	q.Location = dus.queryLocation
	// Run the query and print results when the query job is completed.
	job, err := q.Run(ctx)
	if err != nil {
		return err
	}
	status, err := job.Wait(ctx)
	if err != nil {
		return err
	}
	if err := status.Err(); err != nil {
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
		if iErr == iterator.Done {
			break
		}
		if iErr != nil {
			return iErr
		}
		log.Info().Interface("row", row).Msg("new row")
	}
	return nil
}
