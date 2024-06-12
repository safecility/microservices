package messages

import "github.com/safecility/go/lib"

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

type MilesiteCTReading struct {
	*lib.Device
	UID   string
	Power bool
	Version
	Current
}

type Alarms struct {
	t  bool
	tr bool
	r  bool
	rr bool
}
