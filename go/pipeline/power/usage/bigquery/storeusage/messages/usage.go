package messages

import (
	"cloud.google.com/go/bigquery"
	"github.com/safecility/go/lib"
	"time"
)

type MeterReading struct {
	*lib.Device
	Time       time.Time
	ReadingKWH float64
}

func BigQueryTableMetadata(name string) *bigquery.TableMetadata {
	sampleSchema := bigquery.Schema{
		{Name: "DeviceUID", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "ReadingKWH", Type: bigquery.FloatFieldType},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: sampleSchema,
	}
}
