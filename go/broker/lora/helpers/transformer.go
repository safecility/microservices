package helpers

import (
	"fmt"
	"github.com/safecility/go/mqtt/messages"
)

type SimpleDaliPayloadAdjuster struct{}

func (i SimpleDaliPayloadAdjuster) AdjustPayload(message *messages.LoraMessage) error {
	schema2Payload := append([]byte{0}, message.Payload...)
	message.Payload = schema2Payload
	return nil
}

type AppIdUidTransformer struct {
	AppID string
}

func (u AppIdUidTransformer) GetUID(deviceID string) string {
	return fmt.Sprintf("%s/%s", u.AppID, deviceID)
}

// Segments2Scheme2Adjuster add scheme number, DeviceEUI and 0 messageMode (unused) so payload conforms to schema2
type Segments2Scheme2Adjuster struct{}

func (i Segments2Scheme2Adjuster) AdjustPayload(message *messages.LoraMessage) error {
	schema2Payload := append([]byte{2}, message.DeviceEUI...)
	// the messageMode byte is redundant in LoRA (generally redundant but whatever...)
	schema2Payload = append(schema2Payload, 0)
	schema2Payload = append(schema2Payload, message.Payload...)
	message.Payload = schema2Payload
	return nil
}
