package messages

import (
	"time"
)

type HotdropReading struct {
	DeviceEUI                   string
	MaximumCurrent              float64
	MinimumCurrent              float64
	InstantaneousCurrent        float64
	AverageCurrent              float64
	AccumulatedCurrent          float64
	SecondsAgoForMinimumCurrent float64
	SecondsAgoForMaximumCurrent float64
	SupplyVoltage               float64
	Temp                        float64
}

type HotdropDeviceReading struct {
	*PowerDevice `datastore:",omitempty"`
	HotdropReading
	Time time.Time
}
