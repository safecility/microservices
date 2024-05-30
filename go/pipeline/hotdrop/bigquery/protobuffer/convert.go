package protobuffer

import (
	"github.com/safecility/microservices/go/pipeline/hotdrop/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.HotdropDeviceReading) *Hotdrop {
	return &Hotdrop{
		DeviceEUI:                   r.DeviceEUI,
		Time:                        r.Time.Format(time.RFC3339),
		Temp:                        r.Temp,
		InstantaneousCurrent:        r.InstantaneousCurrent,
		MaximumCurrent:              r.MaximumCurrent,
		SecondsAgoForMaximumCurrent: r.SecondsAgoForMaximumCurrent,
		MinimumCurrent:              r.MinimumCurrent,
		SecondsAgoForMinimumCurrent: r.SecondsAgoForMinimumCurrent,
		AccumulatedCurrent:          r.AccumulatedCurrent,
		SupplyVoltage:               r.SupplyVoltage,
	}
}
