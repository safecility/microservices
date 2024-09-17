package messages

import (
	"github.com/safecility/go/lib/gbigquery"
	"time"
)

type AccumulatorBucket struct {
	Accumulator    string
	AccumulatorUID string
	Max            float64
	Min            float64
	Usage          float64
	BucketStart    time.Time
	BucketInterval gbigquery.TimeInterval
}
