package messages

import (
	"time"
)

type Listing struct {
	SystemUID string `datastore:",omitempty"`
	TenantUID string `datastore:",omitempty"`
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

	Listing    *Listing   `datastore:",flatten,omitempty"`
	Processors *Processor `datastore:",omitempty"`
}

type Processor struct {
	Phase struct {
		PhaseName string `datastore:",omitempty"`
		PhaseUID  string `datastore:",omitempty"`
	}
	System struct {
		SystemUID  string `datastore:",omitempty"`
		SystemName string `datastore:",omitempty"`
		SystemType string `datastore:",omitempty"`
	}
}
