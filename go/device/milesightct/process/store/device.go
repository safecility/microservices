package store

import (
	"github.com/safecility/microservices/go/device/milesightct/process/messages"
)

type DeviceStore interface {
	GetDevice(uid string) (*messages.PowerDevice, error)
}
