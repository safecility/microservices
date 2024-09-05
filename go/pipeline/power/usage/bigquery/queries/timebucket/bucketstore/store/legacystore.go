package store

import (
	"cloud.google.com/go/datastore"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/safecility/go/lib/gbigquery"
	"github.com/safecility/microservices/go/pipeline/power/usage/bigquery/queries/timebucket/bucketstore/messages"
)

type BucketStore interface {
	AddBuckets(m []messages.UsageBucket, ty *gbigquery.BucketType) error
}

type LegacyDatastore struct {
	client *datastore.Client
}

func NewLegacyDatastore(client *datastore.Client) *LegacyDatastore {
	return &LegacyDatastore{client: client}
}

func (d *LegacyDatastore) AddBuckets(m []messages.UsageBucket, ty *gbigquery.BucketType) error {
	if m == nil || len(m) == 0 {
		return fmt.Errorf("empty usageBuckets")
	}
	ctx := context.Background()

	keys := make([]*datastore.Key, len(m))
	for i, b := range m {
		k, err := d.GetBucketKey(&b, ty)
		if err != nil {
			return err
		}
		keys[i] = k
	}

	ret, err := d.client.PutMulti(ctx, keys, m)
	if err != nil {
		return err
	}
	log.Debug().Str("put keys", fmt.Sprintf("%#v", ret)).Str("sample", fmt.Sprintf("%+v", ret[0])).Msg("success")

	return nil
}

func (d *LegacyDatastore) GetBucketKey(u *messages.UsageBucket, i *gbigquery.BucketType) (*datastore.Key, error) {
	var intervalKey string

	switch i.Interval {
	case gbigquery.HOUR:
		intervalKey = fmt.Sprintf("%d-%d-%d %02d", u.Bucket.StartTime.Year(), u.Bucket.StartTime.Month(), u.Bucket.StartTime.Day(), u.Bucket.StartTime.Hour())
	case gbigquery.DAY:
		intervalKey = fmt.Sprintf("%d:d%02d", u.Bucket.StartTime.Year(), u.Bucket.StartTime.YearDay())
	}

	name := fmt.Sprintf("%s_%s", u.DeviceUID, intervalKey)

	return datastore.NameKey("BucketReading", name, nil), nil
}
