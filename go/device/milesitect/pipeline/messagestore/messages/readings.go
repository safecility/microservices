package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

type PowerDevice struct {
	*lib.Device
	PowerFactor float64 `datastore:",omitempty"`
	Voltage     float64 `datastore:",omitempty"`
}

type MilesiteCTReading struct {
	*PowerDevice
	UID   string
	Power bool
	Time  time.Time
	Version
	Current
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}

type Version struct {
	Ipso     string
	Hardware string
	Firmware string
}

type Current struct {
	Total float32
	Value float32
	Max   float32
	Min   float32
	Alarms
}
