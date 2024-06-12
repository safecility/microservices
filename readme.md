### Microservices for streaming device data

Microservices are divided into brokers, processes and pipelines
* Brokers ingest data from devices
* Processes store decoded data and pass messages onwards for further processing in pipelines
* Pipelines create additional data from the Processes interpreted device data

If data is in a known standardized format all is fine (it's usually not) otherwise the ingestion of the maker/device 
specific format is handled by

### Devices
We separate the flow of device specific data

Only once data changes to a common format do the general pipelines take over.

An individual device will typically have its own structure

- device name
- - process
- - pipeline
- - - storage
- - - bigquery
- - - standardize

so a single process microservice to handle raw payload

then standard pipelines to handle storing the data, adding to bigquery, output in some standard format
