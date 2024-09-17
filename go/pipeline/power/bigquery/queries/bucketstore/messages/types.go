package messages

import (
	"github.com/safecility/go/lib/gbigquery"
	"time"
)

const BigTableTime = "2006-01-02 15:04:05"

type UsageBucket struct {
	//DeviceUID, min as initial_reading, max as final_reading, (max - min) as kWh, readings, first_reading, last_reading, bucket
	DeviceUID    string
	Max          float64   `datastore:",omitempty"`
	Min          float64   `datastore:",omitempty"`
	Usage        float64   `datastore:",omitempty"`
	Readings     int       `datastore:",omitempty"`
	FirstReading time.Time `datastore:",omitempty"`
	LastReading  time.Time `datastore:",omitempty"`
	gbigquery.Bucket
	Type string `datastore:",omitempty"`
}
