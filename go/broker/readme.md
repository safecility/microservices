### Data Brokers

Brokers are lightweight pipes between device Data Brokers (NBIoT and LoRA) and our microservices.

The purpose of the broker is to separate Payload processing from the transport.
All brokers create identical messages for processing by microservices

```
type SimpleMessage struct {
	Source     string
	DeviceUID  string
	Payload   []byte
	Time      time.Time
}
```

This separates the Payload processing from the transport handling

To process transport specific information in a microservice create an extension of the SimpleMessage struct

```
type LoraData struct {
	DeviceEUI []byte
	Signal    *Signal
	Channel   MqttChannel
}

type LoraMessage struct {
	SimpleMessage
	LoraData
}
```

now microservices can access the SimpleMessage as usual 
and a microservice with knowledge of the specific transport mechanism can parse the necessary fields
