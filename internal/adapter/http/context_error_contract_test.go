package httpadapter

import (
	"context"
	"errors"
	"github.com/example/telemetry-rule-service/internal/adapter/memory"
	"github.com/example/telemetry-rule-service/internal/application/ingest"
	"testing"
)

func TestClassifyCancellationAndDeadlineCauses(t *testing.T) {
	for _, tc := range []struct {
		name   string
		err    error
		code   string
		status int
	}{
		{"canceled", context.Canceled, "cancelled", 499},
		{"deadline", context.DeadlineExceeded, "deadline_exceeded", 504},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := classify(tc.err)
			if got.Code != tc.code || got.Status != tc.status {
				t.Fatalf("classification = %s/%d, want %s/%d", got.Code, got.Status, tc.code, tc.status)
			}
		})
	}
}
func TestTelemetryIngestErrorClassificationAcrossLayers(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := ingest.NewService(memory.MetricRepository{Store: memory.NewStore()}).IngestJSON(ctx, []byte(`{"source_id":"plant-a","metric_id":"temperature","value":91}`))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("JSON cancellation lost across layers: %v", err)
	}
	got := classify(err)
	if got.Code != "cancelled" || got.Status != 499 {
		t.Fatalf("JSON cancellation classified as %s/%d: %v", got.Code, got.Status, err)
	}
	if got.Message != "request canceled: record sample: record metric sample: context canceled" {
		t.Fatalf("JSON cancellation message = %q", got.Message)
	}
}
