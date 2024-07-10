# Create System Buckets 

Timebuckets based on the SystemUID - all devices in a time period which are part of a system (share the same SystemUID) 
have their kWh usage summed 

`````sql
SELECT SystemUID, SUM(max - min) as kWh, COUNT(DeviceUID), bucket from (
  SELECT SystemUID, DeviceUID, Max(ReadingKWH) as max, Min(ReadingKWH) as min, TIMESTAMP_BUCKET(Time, INTERVAL 1 HOUR) AS bucket 
  FROM **bigquery-table**
  WHERE SystemUID != ""
  GROUP BY SystemUID, DeviceUID, bucket
)
GROUP BY SystemUID, bucket
ORDER BY bucket;
`````