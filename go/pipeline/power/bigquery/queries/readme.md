## Query servers

These servers run queries on stored bigquery values

### Accumulator

Accumulator runs on a device's accumulator values, for example
a Device with a 'phase' accumulator. All devices belonging to a phase
will have their values added and the output added to the phase Topic

### Bucketstore

This runs a server that creates values from time windows - hour or day

The time window is stored in datastore for each device