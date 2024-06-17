package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usagestore/messages"
)

type DatastoreUsage struct {
	client *datastore.Client
}

func NewDatastoreUsage(client *datastore.Client) (*DatastoreUsage, error) {
	rd := &DatastoreUsage{client: client}
	return rd, nil
}

func (d *DatastoreUsage) AddMeterReading(m *messages.MeterReading) error {
	ctx := context.Background()
	k := datastore.IncompleteKey("MeterReading", nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}
	log.Debug().Str("uid", m.DeviceUID).Msg("putting new meter reading")
	return nil
}
