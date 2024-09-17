## Store Three Phase usage
devices in a 3 phase system have their usage stored in bigquery

periodic queries output summed usage for given time periods

this microservice takes messages from the phase topic and records them

messages have a AccumulatorUID identifier - from this we retrieve a virtual Output Device
which acts as a standard device