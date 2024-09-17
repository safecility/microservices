package store

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/safecility/go/lib"
	"time"
)

type DeviceFirestore struct {
	client          *firestore.Client
	contextDeadline time.Duration
}

func NewDeviceFirestore(client *firestore.Client, deadline int) *DeviceFirestore {
	return &DeviceFirestore{client: client, contextDeadline: time.Duration(deadline)}
}

func (df DeviceFirestore) GetDevice(uid string) (*lib.Device, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*df.contextDeadline)
	defer cancel()

	m, err := df.client.Collection("device").Doc(uid).Get(ctx)
	if err != nil {
		return nil, err
	}
	d := &lib.Device{
		DeviceMeta: &lib.DeviceMeta{
			Firmware:   &lib.Firmware{},
			Processors: &lib.Processor{},
		},
	}
	err = m.DataTo(d)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (df DeviceFirestore) Close() error {
	return df.client.Close()
}
