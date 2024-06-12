package messages

import "github.com/safecility/go/lib"

type EastronEagleReading struct {
	*lib.Device
	ImportActiveEnergy   float64
	ExportActiveEnergy   float64
	ActivePower          float64
	InstantaneousCurrent float64
	InstantaneousVoltage float64
	PowerFactor          float64
	RelayState           float64
}
