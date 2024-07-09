package messages

import (
	"github.com/safecility/go/lib"
)

// PowerProfile - not all our power meters have access to instantaneous PowerFactor and Voltage information -
// to generate kWh readings they use an average stored value to convert from Amp Hour to KWh - these values are cached
// with the device information and passed down the pipeline
type PowerProfile struct {
	PowerFactor float64 `firestore:",omitempty"`
	Voltage     float64 `firestore:",omitempty"`
}

// PowerDevice We use datastore for Device message processing so there is no need to store PowerProfile information
type PowerDevice struct {
	lib.Device
	PowerProfile *PowerProfile `datastore:"-" firestore:",omitempty"`
}
