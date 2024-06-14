package messages

import (
	"fmt"
	"github.com/safecility/go/lib"
	"time"
)

type PowerDevice struct {
	*lib.Device
	PowerFactor float64 `datastore:",omitempty"`
	Voltage     float64 `datastore:",omitempty"`
}

type MilesightCTReading struct {
	*PowerDevice
	UID   string
	Power bool
	Time  time.Time
	Version
	Current
}

type MeterReading struct {
	*lib.Device
	ReadingKWH float64
	Time       time.Time
}

func (mc MilesightCTReading) Usage() (*MeterReading, error) {
	if mc.PowerDevice == nil {
		return nil, fmt.Errorf("device does not have its PowerDevice definitions")
	}
	kWh := float64(mc.Current.Total) * mc.Voltage * mc.PowerFactor
	return &MeterReading{
		Device:     mc.PowerDevice.Device,
		ReadingKWH: kWh,
		Time:       mc.Time,
	}, nil
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