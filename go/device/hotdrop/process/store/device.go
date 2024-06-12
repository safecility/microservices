package store

import (
	"github.com/safecility/microservices/go/device/hotdrop/process/messages"
)

type DeviceStore interface {
	GetDevice(uid string) (*messages.PowerDevice, error)
}
