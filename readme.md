# Microservices

Microservices for operating on IoT device data

Currently all microservices are golang see the contained directory for details

## Structure

* Transports: handle messages from data-brokers and webhook input
* Process: convert payloads to storage and pubsub messages (these are all stored in device repos)
* Pipelines: device repo's have their own pipelines - 
pipelines contained here are for general message types output from device pipelines