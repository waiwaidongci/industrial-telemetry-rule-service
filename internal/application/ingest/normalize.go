package ingest

import (
	"context"
	"fmt"
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"strings"
	"time"
)

type SourceLookup interface {
	Get(context.Context, string) (interface{}, error)
}
type Validator struct {
	MaxAge    time.Duration
	MaxFuture time.Duration
}

func (v Validator) Validate(sample metric.Sample, now time.Time) error {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if err := sample.ValidateStrict(); err != nil {
		return err
	}
	if v.MaxAge > 0 && now.Sub(sample.Timestamp) > v.MaxAge {
		return fmt.Errorf("sample is older than %s", v.MaxAge)
	}
	if v.MaxFuture > 0 && sample.Timestamp.After(now.Add(v.MaxFuture)) {
		return fmt.Errorf("sample is from the future")
	}
	return nil
}
func Normalize(sample metric.Sample) metric.Sample {
	sample.SourceID = strings.TrimSpace(sample.SourceID)
	sample.MetricID = strings.TrimSpace(sample.MetricID)
	normalized := make(map[string]string, len(sample.Tags))
	for key, value := range sample.Tags {
		key = strings.ToLower(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		if key != "" && value != "" {
			normalized[key] = value
		}
	}
	clear(sample.Tags)
	for key, value := range normalized {
		sample.Tags[key] = value
	}
	if sample.Timestamp.IsZero() {
		sample.Timestamp = time.Now().UTC()
	}
	return sample
}
