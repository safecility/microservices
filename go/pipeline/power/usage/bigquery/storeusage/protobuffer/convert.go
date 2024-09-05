package protobuffer

import (
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/store/messages"
	"time"
)

var (
	threePhase = "threePhase"
	solar      = "solar"
)

func CreateProtobufMessage(r *messages.MeterReading) *Usage {

	usage := &Usage{
		DeviceUID:  r.DeviceUID,
		Time:       r.Time.Format(time.RFC3339),
		ReadingKWH: r.ReadingKWH,

		DeviceName:  r.DeviceName,
		DeviceTag:   r.DeviceTag,
		CompanyUID:  r.CompanyUID,
		LocationUID: r.LocationUID,
	}

	if r.Processors != nil {
		for k, p := range *r.Processors {
			var uid string
			var name string
			processor, ok := p.(map[string]interface{})
			if ok {
				u := processor["uid"]
				if u != nil {
					uid = u.(string)
				} else {
					log.Warn().Msg("error parsing processor UID")
				}
				n := processor["name"]
				if u != nil {
					name = n.(string)
				} else {
					log.Warn().Msg("error parsing processor name")
				}
			}

			switch k {
			case "solar":
				usage.SystemUID = &uid
				usage.SystemName = &name
				usage.SystemType = &solar
				el := processor["element"]
				if el != nil {
					se := el.(string)
					usage.SystemElement = &se
				} else {
					log.Warn().Msg("error parsing processor name")
				}
				break
			case "phase":
				usage.SystemUID = &uid
				usage.SystemName = &name
				usage.SystemType = &threePhase
				p := processor["phase"]
				if p != nil {
					se := p.(string)
					usage.SystemElement = &se
				} else {
					log.Warn().Msg("error parsing processor name")
				}
				break
			case "tenant":
				usage.TenantUID = &uid
				usage.TenantName = &name
				break
			case "group":
				usage.GroupUID = &uid
				usage.GroupName = &name
				break
			}
		}
	}

	return usage
}
