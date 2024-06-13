package protobuffer

import (
	"github.com/safecility/microservices/go/device/eastronsdm/pipeline/bigquery/messages"
	"time"
)

func CreateProtobufMessage(r *messages.EastronSdmReading) *EastronSdmBq {
	return &EastronSdmBq{
		DeviceUID:            r.DeviceUID,
		Time:                 r.Time.Format(time.RFC3339),
		ActivePower:          r.ActivePower,
		ImportActiveEnergy:   r.ImportActiveEnergy,
		ExportActiveEnergy:   r.ExportActiveEnergy,
		InstantaneousCurrent: r.InstantaneousCurrent,
		InstantaneousVoltage: r.InstantaneousVoltage,
		PowerFactor:          r.PowerFactor,
	}
}
