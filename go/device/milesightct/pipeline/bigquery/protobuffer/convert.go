package protobuffer

import (
	"github.com/safecility/microservices/go/device/milesightct/pipeline/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.MilesightCTReading) *Milesight {
	return &Milesight{
		DeviceUID:            r.DeviceUID,
		Time:                 r.Time.Format(time.RFC3339),
		AccumulatedCurrent:   float64(r.Current.Total),
		InstantaneousCurrent: float64(r.Current.Value),
		MaximumCurrent:       float64(r.Current.Max),
		MinimumCurrent:       float64(r.Current.Min),
	}
}
