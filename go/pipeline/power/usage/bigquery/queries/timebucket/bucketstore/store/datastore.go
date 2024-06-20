package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket"
)

type DatastoreBuckets struct {
	client *datastore.Client
}

func NewBucketDatastore(client *datastore.Client) *DatastoreBuckets {
	return &DatastoreBuckets{client: client}
}

func (d *DatastoreBuckets) AddBucket(m *timebucket.UsageBucket, i *timebucket.BucketType) error {
	ctx := context.Background()
	k, err := d.GetBucketKey(m, i)
	if err != nil {
		return err
	}
	k, err = d.client.Put(ctx, k, m)
	if err != nil {
		return err
	}
	log.Debug().Str("uid", m.DeviceUID).Msg("put new bucket")
	return nil
}

func (d *DatastoreBuckets) AddBuckets(m []timebucket.UsageBucket, i *timebucket.BucketType) error {
	ctx := context.Background()

	keys := make([]*datastore.Key, len(m))
	for _, b := range m {
		k, err := d.GetBucketKey(&b, i)
		if err != nil {
			return err
		}
		keys = append(keys, k)
	}

	ret, err := d.client.PutMulti(ctx, keys, m)
	if err != nil {
		return err
	}
	log.Debug().Str("put keys", fmt.Sprintf("%#v", ret)).Msg("success")

	return nil
}

func (d *DatastoreBuckets) GetBucketKey(u *timebucket.UsageBucket, i *timebucket.BucketType) (*datastore.Key, error) {
	var intervalKey string

	switch i.Interval {
	case timebucket.HOUR:
		intervalKey = fmt.Sprintf("%d-%d-%d %02d", u.Bucket.StartTime.Year(), u.Bucket.StartTime.Month(), u.Bucket.StartTime.Day(), u.Bucket.StartTime.Hour())
	case timebucket.DAY:
		intervalKey = fmt.Sprintf("%d:d%02d", u.Bucket.StartTime.Year(), u.Bucket.StartTime.YearDay())
	}

	name := fmt.Sprintf("%s_%s", u.DeviceUID, intervalKey)

	return datastore.NameKey("MeterReading", name, nil), nil
}
