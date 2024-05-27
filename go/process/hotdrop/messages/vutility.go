package messages

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib"
	"strings"
	"time"
)

type HotdropReading struct {
	DeviceEUI      string
	MaximumCurrent float64
	MinimumCurrent float64
	//RMS average rate of current instantaneousCurrent_MilliAmpere (for LoRA in nano_Amps?)
	InstantaneousCurrent float64
	AverageCurrent       float64
	//Nano Ampere HourAmps
	AccumulatedCurrent          float64
	SecondsAgoForMinimumCurrent float64
	SecondsAgoForMaximumCurrent float64
	SupplyVoltage               float64
	Temp                        float64
}

type HotdropDeviceReading struct {
	lib.Device
	HotdropReading
	Time time.Time
}

type VuSensorMessage struct {
	Data []struct {
		DevEui                    string    `json:"devEui"`
		ApiReceivedAt             time.Time `json:"apiReceivedAt"`
		ExternalNetworkType       string    `json:"externalNetworkType"`
		ExternalNetworkName       string    `json:"externalNetworkName"`
		ExternalNetworkReceivedAt time.Time `json:"externalNetworkReceivedAt"`
		Rssi                      float64   `json:"rssi"`
		Snr                       float64   `json:"snr"`
		FrameCount                int       `json:"frameCount"`
		Latitude                  float64   `json:"latitude"`
		Longitude                 float64   `json:"longitude"`
		Altitude                  float64   `json:"altitude"`
		SensorMeasurements        []struct {
			Type  string  `json:"type"`
			Value float64 `json:"value"`
		} `json:"sensorMeasurements"`
	} `json:"data"`
}

func (m *VuSensorMessage) GetHotDropReadings() (data []HotdropDeviceReading) {

	dataMap := make(map[string]int, len(m.Data))

	for i, d := range m.Data {

		measurements := len(d.SensorMeasurements)

		if measurements == 0 {
			log.Warn().Interface("data", d).Msg("no sensor measurements")
			continue
		}

		key := fmt.Sprintf("%s:%d", d.DevEui, d.ExternalNetworkReceivedAt.Nanosecond())
		_, e := dataMap[key]
		if e {
			log.Warn().Int("count", i).
				Interface("vutility", m).Msg("vutility message contains duplicates")
			continue
		}
		eui := strings.ToLower(d.DevEui)

		hdr := HotdropDeviceReading{
			Device: lib.Device{
				DeviceUID: eui,
			},
			Time: d.ExternalNetworkReceivedAt,
			HotdropReading: HotdropReading{
				DeviceEUI: eui,
			},
		}
		dataMap[key] = 1

		// The api is unfortunately unstable - we use milli values as these have always been available
		// Similarly accumulatedCurrent_NanoAmpereHour and / 1000000 - annoying
		for _, hdm := range d.SensorMeasurements {
			switch hdm.Type {
			case "temperature_Celsius":
				hdr.Temp = hdm.Value
				break
			case "instantaneousCurrent_MilliAmpere":
				hdr.InstantaneousCurrent = hdm.Value
				break
			case "maximumCurrent_MilliAmpere":
				hdr.MaximumCurrent = hdm.Value
				break
			case "maximumCurrent_Ampere":
				break
			case "secondsAgoForMaximumCurrent":
				hdr.SecondsAgoForMinimumCurrent = hdm.Value
				break
			case "minimumCurrent_MilliAmpere":
				hdr.MinimumCurrent = hdm.Value
				break
			case "minimumCurrent_Ampere":
				break
			case "secondsAgoForMinimumCurrent":
				hdr.SecondsAgoForMinimumCurrent = hdm.Value
				break
			case "accumulatedCurrent_NanoAmpereHour":
				hdr.AccumulatedCurrent = hdm.Value / 1000000
				break
			case "accumulatedCurrent_AmpereHour":
				break
			case "supplyVoltage_Volt":
				hdr.SupplyVoltage = hdm.Value
				break
			case "averageCurrent_Ampere":
				hdr.AverageCurrent = hdm.Value
				break
			case "resistorIndex":
				break
			default:
				log.Warn().Str("type", hdm.Type).Msg("unknown measurement")
			}
		}
		if hdr.AccumulatedCurrent == 0 {
			log.Warn().Interface("HotDrop", hdr).Msg("Accumulated current 0")
			continue
		}
		data = append(data, hdr)
	}

	return data
}
