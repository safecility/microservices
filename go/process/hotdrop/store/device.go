package store

import (
	"github.com/safecility/microservices/go/process/hotdrop/messages"
)

type DeviceStore interface {
	GetDevice(uid string) (*messages.PowerDevice, error)
}
