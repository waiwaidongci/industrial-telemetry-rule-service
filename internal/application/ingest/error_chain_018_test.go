package ingest

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

var errRecord018 = errors.New("record sentinel 018")

type recordFailure018 struct{}

func (recordFailure018) Record(context.Context, metric.Sample) error { return errRecord018 }

func sample018() metric.Sample {
	return metric.Sample{SourceID: "source", MetricID: "metric", Value: 1, Timestamp: time.Now().UTC()}
}

func TestBatchPreservesRecordCause018(t *testing.T) {
	_, err := NewService(recordFailure018{}).IngestBatch(context.Background(), []metric.Sample{sample018()})
	if !errors.Is(err, errRecord018) { t.Fatalf("batch lost cause: %v", err) }
}

func TestPipelinePreservesRecordCause018(t *testing.T) {
	p := NewPipeline(NewService(recordFailure018{}), Validator{}, nil)
	if err := p.Process(context.Background(), sample018()); !errors.Is(err, errRecord018) { t.Fatalf("pipeline lost cause: %v", err) }
}
