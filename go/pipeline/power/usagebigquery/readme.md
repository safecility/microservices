### Publish Bigquery compatible messages from Power Usage

To store data in Bigquery via pubsub we need 
* a bigQuery table with a schema
* a pubsub.Topic using the same schema
* a pubsub.Subscription set to handle bigQuery

(see: https://cloud.google.com/pubsub/docs/create-bigquery-subscription)