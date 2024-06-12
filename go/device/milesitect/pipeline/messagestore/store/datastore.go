package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/device/milesitect/pipeline/messagestore/messages"
)

type DatastoreMilesite struct {
	client *datastore.Client
}

func NewDatastoreMilesite(client *datastore.Client) (*DatastoreMilesite, error) {
	rd := &DatastoreMilesite{client: client}
	return rd, nil
}

func (d *DatastoreMilesite) AddMilesiteMessage(m *messages.MilesiteCTReading) error {
	ctx := context.Background()
	k := datastore.IncompleteKey("MilesiteCT", nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}
	log.Debug().Str("uid", m.DeviceUID).Msg("putting new eastron eagle message")
	return nil
}
