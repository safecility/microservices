### Power usage

simple microservices built on streams of MeterReadings
time windows are hourly, daily, weekly, monthly

we work from a cache - so the first reading of any period is stored in cache. 
We output a UsagePeriod message with the first measurement of the current period and current reading.

If no first measurement exists the current reading is added to the cache as its first reading.


