## Webhooks

Provide transports that deal with known webhook formats

Both the Vutility and Everynet APIs wrap a LoRA delivery system: At the moment we concentrate on Uplink messages but we
will shortly deliver downlink, receipts etc as with the mqtt (lora/paho) transport

The webhook formats are general - we could configure the webhook or the jwt to specify a device type and route uplink
messages to a specific topic based on the device type
* jwt topic configuration solution has the advantage that additional topics can be configured outside the microservice

### Everynet
everynet provides LoRA message delivery through a webhook api

### Vutility
vutility provides a generalized sensor message -> we provide a process microservice to interpret the data from Vutility's
hotdrop device -> as the Vutility API is only used by us for the Hotdrop device its webhook has moved to the hotdrop device
repo https://github.com/safecility/hotdrop/
