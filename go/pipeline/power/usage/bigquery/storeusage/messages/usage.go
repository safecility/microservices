package messages

import (
	"cloud.google.com/go/bigquery"
	"github.com/safecility/go/lib"
	"time"
)

type PowerProfile struct {
	PowerFactor float64 `firestore:",omitempty" json:"powerFactor,omitempty"`
	Voltage     float64 `firestore:",omitempty" json:"voltage,omitempty"`
}

type PowerDevice struct {
	lib.Device
	PowerProfile *PowerProfile `datastore:"-" firestore:",omitempty" json:"powerProfile,omitempty"`
}

type MeterReading struct {
	PowerDevice
	Time       time.Time
	ReadingKWH float64
}

// BigQueryTableMetadata TODO can we get this from the protobuf file?
func BigQueryTableMetadata(name string) *bigquery.TableMetadata {

	usageSchema := bigquery.Schema{
		{Name: "DeviceUID", Type: bigquery.StringFieldType},
		{Name: "Time", Type: bigquery.TimestampFieldType},
		{Name: "ReadingKWH", Type: bigquery.FloatFieldType},

		{Name: "DeviceName", Type: bigquery.StringFieldType},
		{Name: "DeviceTag", Type: bigquery.StringFieldType},
		{Name: "CompanyUID", Type: bigquery.StringFieldType},
		{Name: "LocationUID", Type: bigquery.StringFieldType},

		{Name: "SystemUID", Type: bigquery.StringFieldType, Required: false},
		{Name: "TenantUID", Type: bigquery.StringFieldType, Required: false},
		{Name: "GroupUID", Type: bigquery.StringFieldType, Required: false},
	}

	return &bigquery.TableMetadata{
		Name:   name,
		Schema: usageSchema,
	}
}
