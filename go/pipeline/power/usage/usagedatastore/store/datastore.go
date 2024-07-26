package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usage/store/messages"
)

type DatastoreUsage struct {
	client *datastore.Client
	entity string
}

func NewDatastoreUsage(client *datastore.Client, entity string) (*DatastoreUsage, error) {
	rd := &DatastoreUsage{client: client, entity: entity}
	return rd, nil
}

func (d *DatastoreUsage) AddMeterReading(m *messages.MeterReading) error {
	ctx := context.Background()
	k := datastore.IncompleteKey(d.entity, nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}

	log.Debug().Str("entity", d.entity).Str("uid", m.DeviceUID).Msg("put new meter reading")
	return nil
}
