package messages

import "github.com/safecility/go/lib/gbigquery"

type UsageBucket struct {
	SystemUID string
	Usage     float64 `datastore:",omitempty"`
	gbigquery.Bucket
	Type string `datastore:",omitempty"`
}
