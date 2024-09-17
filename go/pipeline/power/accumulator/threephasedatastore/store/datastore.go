package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/accumulator/threephase/messages"
)

type ThreePhaseDatastore struct {
	client *datastore.Client
	entity string
}

func NewThreePhaseDatastore(client *datastore.Client, entity string) *ThreePhaseDatastore {
	return &ThreePhaseDatastore{client: client, entity: entity}
}

func (d *ThreePhaseDatastore) AddDevicePowerUse(m *messages.MeterReading) error {

	ctx := context.Background()

	n := d.readingKey(m)
	k := datastore.NameKey(d.entity, n, nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}

	log.Debug().Str("entity", d.entity).Str("uid", m.DeviceUID).Msg("put new meter reading")
	return nil
}

func (d *ThreePhaseDatastore) AddMeterReading(m *messages.MeterReading) error {

	ctx := context.Background()

	n := d.readingKey(m)
	k := datastore.NameKey(d.entity, n, nil)
	k, err := d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}

	log.Debug().Str("entity", d.entity).Str("uid", m.DeviceUID).Msg("put new meter reading")
	return nil
}

func (d *ThreePhaseDatastore) readingKey(m *messages.MeterReading) string {
	k := m.DeviceUID
	t := m.Time
	switch m.Interval {
	case gbigquery.DAY:
		k = fmt.Sprintf("%s:d-%s", k, t.Format("2006-01-02"))
	case gbigquery.HOUR:
		k = fmt.Sprintf("%s:h-%s", k, t.Format("2006-01-02T15"))
	}
	return k
}
