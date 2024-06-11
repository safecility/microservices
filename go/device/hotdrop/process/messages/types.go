package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

type PowerReading struct {
	Reading float64
	Units   string
	Time    time.Time
}

type PowerDevice struct {
	*lib.Device
	PowerFactor float64 `datastore:",omitempty"`
	Voltage     float64 `datastore:",omitempty"`
}
