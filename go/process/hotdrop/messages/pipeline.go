package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

// PowerUsage removes Hotdrop formatting and by using Unit should be unambiguously processable by other services
// because we use power factor and voltage here they don't need to be passed onwards in the pipeline
type PowerUsage struct {
	*lib.Device
	Reading float64
	Units   string
	Time    time.Time
}

func GetPowerUsage(hd HotdropDeviceReading, device *PowerDevice) *PowerUsage {
	reading := hd.AccumulatedCurrent * device.Voltage * device.PowerFactor
	return &PowerUsage{Device: device.Device, Reading: reading, Units: "kWh", Time: hd.Time}
}
