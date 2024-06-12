package protobuffer

import (
	"github.com/safecility/microservices/go/device/eagle/pipeline/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.EastronEagleReading) *Hotdrop {
	return &Hotdrop{
		DeviceEUI:                   r.DeviceEUI,
		Time:                        r.Time.Format(time.RFC3339),
		Temp:                        r.Temp,
		InstantaneousCurrent:        r.InstantaneousCurrent.Milli,
		MaximumCurrent:              r.MaximumCurrent.Milli,
		SecondsAgoForMaximumCurrent: r.SecondsAgoForMaximumCurrent,
		MinimumCurrent:              r.MinimumCurrent.Milli,
		SecondsAgoForMinimumCurrent: r.SecondsAgoForMinimumCurrent,
		AccumulatedCurrent:          r.AccumulatedCurrent.Milli,
		SupplyVoltage:               r.SupplyVoltage,
	}
}