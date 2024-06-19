## Pipelines
These effectively represent *all the rest* of the microservices

### Non specific
Pipelines where messages have a generalized form appear here - otherwise they will be present under the device

### Organization

Different device classes have their own trees and message types within those have their own trees.
    
* Power
  * Usage
    * Datastore
    * Bigquery
      * Store
      * Queries...

In power we have a pipeline for Usage (meter readings) we store these in Datastore and Bigquery

There are additional bigquery microservices that run queries on the stored data
   
        