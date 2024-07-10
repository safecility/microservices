package protobuffer

import (
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/messages"
	"time"
)

func CreateProtobufMessage(r *messages.MeterReading) *Usage {
	return &Usage{
		DeviceUID:  r.DeviceUID,
		Time:       r.Time.Format(time.RFC3339),
		ReadingKWH: r.ReadingKWH,

		DeviceName:  r.DeviceName,
		DeviceTag:   r.DeviceTag,
		CompanyUID:  r.CompanyUID,
		LocationUID: r.LocationUID,

		SystemUID: r.Listing.SystemUID,
		TenantUID: r.Listing.TenantUID,
	}
}
