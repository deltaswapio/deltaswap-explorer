import "date"

option task = {
    name: "chain activity for 30 days with 3-hour granularity",
    every: 3h,
}

sourceBucket = "deltaswapscan"
destinationBucket = "deltaswapscan-24hours"
execution = date.truncate(t: now(), unit: 1h)
start = date.truncate(t: -30d, unit: 24h)

from(bucket: sourceBucket)
  |> range(start: start)
  |> filter(fn: (r) => r._measurement == "vaa_volume_v2" and r._field == "volume")
  |> group(columns: ["emitter_chain", "destination_chain", "app_id"])
  |> count(column: "_value")
  |> map(fn: (r) => ({r with _time: execution}))
  |> set(key: "_measurement", value: "chain_activity_30_days_3h_v2")
  |> set(key: "_field", value: "count")
  |> to(bucket: destinationBucket)

from(bucket: sourceBucket)
  |> range(start: start)
  |> filter(fn: (r) => r._measurement == "vaa_volume_v2" and r._field == "volume")
  |> group(columns: ["emitter_chain", "destination_chain", "app_id"])
  |> sum(column: "_value")
  |> map(fn: (r) => ({r with _time: execution}))
  |> set(key: "_measurement", value: "chain_activity_30_days_3h_v2")
  |> set(key: "_field", value: "notional")
  |> to(bucket: destinationBucket)
