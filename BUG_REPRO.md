# Bug Reproduction

## Bug

Concurrent sample writes, reads, and cleanup trigger the Go race detector. Sample snapshots returned by the memory store also share their `Tags` map with stored or caller-owned data, so tags in an already returned sample can change later.

## Trigger

Run the record's six targeted race-enabled tests against the red branch. The ownership tests mutate the input or returned tag maps after `RecordSample`, `Samples`, `QuerySamples`, and `Latest`. The concurrency test overlaps `RecordSample`, `Samples`, `Latest`, and `ClearSamples` from multiple goroutines.

## Observed Errors

All six targeted tests exit with status 1. The ownership checks report messages including:

```text
clone shares tags: "changed"
stored sample follows caller map: "changed"
Samples leaked mutable tags: "changed"
QuerySamples leaked mutable tags: "changed"
Latest leaked mutable tags: "changed"
```

The concurrent check reports:

```text
WARNING: DATA RACE
```

The race report includes conflicting access in `Store.Samples`, `Store.Latest`, `Store.RecordSample`, and `Store.ClearSamples`.
