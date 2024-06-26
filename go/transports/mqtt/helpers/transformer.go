package helpers

import (
	"fmt"
)

type AppIdUidTransformer struct {
	AppID string
}

func (u AppIdUidTransformer) GetUID(deviceID string) string {
	return fmt.Sprintf("%s/%s", u.AppID, deviceID)
}
