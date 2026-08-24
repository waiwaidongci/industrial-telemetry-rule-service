# Bug Reproduction

## Symptom

Errors returned by batch ingest and protocol decoding lose their original error identity. Callers cannot reliably use `errors.Is` or `errors.As` after the error crosses the batch, pipeline, protocol registry, or service layers.

## Reproduction

From the repository root, run:

```text
go test ./internal/application/ingest -run '^TestBatchPreservesRecordCause018$' -count=1
go test ./internal/application/ingest -run '^TestPipelinePreservesRecordCause018$' -count=1
go test ./internal/adapter/protocol -run '^TestRegistryPreservesDecodeCause018$' -count=1
go test ./internal/adapter/protocol -run '^TestRegistryPreservesUnsupportedCause018$' -count=1
```

On the buggy baseline these checks fail because intermediate layers flatten wrapped errors with `%v` or store them as strings. The corrected implementation preserves the wrapping chain with `%w` and keeps batch errors as `error` values, so all four checks pass.
