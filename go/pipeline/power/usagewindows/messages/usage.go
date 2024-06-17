package messages

import (
	"github.com/safecility/go/lib"
	"time"
)

type MeterReading struct {
	*lib.Device
	ReadingKWH float64
	Time       time.Time
}
