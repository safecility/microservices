## Everynet

Everynet provides a LoRA-webhook IoT solution

The everynet webhook json message double wraps content:
* a type wrapper with a contained RawMessage
* the wrapped content e.g UplinkParams which contains a Payload string

### Pubsub

As with the other webhook microservices currently there is a single webhook to topic linkage.

### JWT Topics

It would be trivial to further support a more generalized service for device types by encoding a pubsub Topic 
within the JWT bearer - at present this is not seen as a priority as the current setup allows good separation 
of services and control over individual device types via gcloud.
