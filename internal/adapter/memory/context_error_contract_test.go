package memory

import (
	"context"
	"errors"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"testing"
	"time"
)

func TestMetricRepositoryRecordPreservesContextCause(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := (MetricRepository{Store: NewStore()}).Record(ctx, metric.Sample{SourceID: "plant-a", MetricID: "temperature", Value: 91, Timestamp: time.Now()})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("metric repository lost cancellation cause: %v", err)
	}
}
