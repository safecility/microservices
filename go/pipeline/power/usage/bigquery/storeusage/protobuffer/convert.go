package protobuffer

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/messages"
	"time"
)

func CreateProtobufMessage(r *messages.MeterReading) *Usage {

	var systemUID, tenantUID, groupUID string
	if r.Processors != nil {
		for k, p := range *r.Processors {
			var uid string
			system, ok := p.(map[string]interface{})
			if ok {
				i := system["SystemUID"]
				if i != nil {
					uid = i.(string)
				} else {
					log.Warn().Msg("error parsing system UID")
				}
			}

			switch k {
			case "Solar":
				systemUID = uid
				break
			case "Phase":
				systemUID = uid
				break
			case "Tenant":
				tenantUID = uid
				break
			case "Group":
				groupUID = uid
				break
			}
		}
	}

	return &Usage{
		DeviceUID:  r.DeviceUID,
		Time:       r.Time.Format(time.RFC3339),
		ReadingKWH: r.ReadingKWH,

		DeviceName:  r.DeviceName,
		DeviceTag:   r.DeviceTag,
		CompanyUID:  r.CompanyUID,
		LocationUID: r.LocationUID,

		SystemUID: systemUID,
		TenantUID: tenantUID,
		GroupUID:  groupUID,
	}
}
