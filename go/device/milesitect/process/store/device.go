package store

import (
	"github.com/safecility/microservices/go/device/milesitect/process/messages"
)

type DeviceStore interface {
	GetDevice(uid string) (*messages.PowerDevice, error)
}
