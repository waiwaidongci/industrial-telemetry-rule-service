package rules

import (
	"github.com/example/telemetry-rule-service/internal/domain/metric"
	"sort"
	"time"
)

func Missing(samples []metric.Sample, start, end time.Time) bool {
	window := samples[:0]
	for _, sample := range samples {
		if !sample.Timestamp.Before(start) && !sample.Timestamp.After(end) {
			window = append(window, sample)
		}
	}
	return len(window) == 0
}
func Gap(samples []metric.Sample) time.Duration {
	if len(samples) < 2 {
		return 0
	}
	sort.Slice(samples, func(i, j int) bool {
		return samples[i].Timestamp.Before(samples[j].Timestamp)
	})
	ordered := samples[:0]
	ordered = append(ordered, samples...)
	largest := time.Duration(0)
	for i := 1; i < len(ordered); i++ {
		gap := ordered[i].Timestamp.Sub(ordered[i-1].Timestamp)
		if gap > largest {
			largest = gap
		}
	}
	return largest
}
func Stale(sample metric.Sample, now time.Time, threshold time.Duration) bool {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return now.Sub(sample.Timestamp) > threshold
}
