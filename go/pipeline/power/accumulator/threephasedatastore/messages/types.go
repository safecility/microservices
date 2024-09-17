package messages

import (
	"github.com/safecility/go/lib/gbigquery"
	"time"
)

type AccumulatorBucket struct {
	Accumulator    string
	AccumulatorUID string
	Max            float64
	Min            float64
	Usage          float64
	BucketStart    time.Time
	BucketInterval gbigquery.TimeInterval
}

type MeterReading struct {
	DeviceUID   string
	DeviceName  string `datastore:",omitempty"`
	DeviceTag   string `datastore:",omitempty"`
	DeviceType  string `datastore:",omitempty"`
	CompanyUID  string `datastore:",omitempty"`
	LocationUID string `datastore:",omitempty"`

	ReadingKWH float64
	Time       time.Time
	Interval   gbigquery.TimeInterval
}
