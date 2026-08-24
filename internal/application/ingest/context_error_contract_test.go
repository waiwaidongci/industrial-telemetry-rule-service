package ingest

import (
	"context"
	"errors"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"testing"
)

type errorRepo struct{ err error }

func (r errorRepo) Record(context.Context, metric.Sample) error { return r.err }
func sampleJSON() []byte {
	return []byte(`{"source_id":"plant-a","metric_id":"temperature","value":91,"timestamp":"2024-01-01T00:00:00Z"}`)
}

func TestIngestJSONPreservesCanceledCause(t *testing.T) {
	_, err := NewService(errorRepo{err: context.Canceled}).IngestJSON(context.Background(), sampleJSON())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("JSON ingest lost cancellation cause: %v", err)
	}
}
func TestIngestTextPreservesDeadlineCause(t *testing.T) {
	_, err := NewService(errorRepo{err: context.DeadlineExceeded}).IngestText(context.Background(), "plant-a,temperature,91,2024-01-01T00:00:00Z")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("text ingest lost deadline cause: %v", err)
	}
}
