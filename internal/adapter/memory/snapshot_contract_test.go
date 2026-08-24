package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

func sampleForSnapshot(index int) metric.Sample {
	return metric.Sample{
		SourceID:  "line-a",
		MetricID:  "temperature",
		Value:     float64(index),
		Timestamp: time.Unix(int64(index+1), 0).UTC(),
		Tags:      map[string]string{"sequence": fmt.Sprint(index)},
	}
}

func TestCloneSampleCopiesTags(t *testing.T) {
	original := sampleForSnapshot(1)
	cloned := cloneSample(original)
	original.Tags["sequence"] = "changed"
	if cloned.Tags["sequence"] != "1" {
		t.Fatalf("clone shares tags: %q", cloned.Tags["sequence"])
	}
}

func TestRecordSampleOwnsInputTags(t *testing.T) {
	store := NewStore()
	value := sampleForSnapshot(2)
	if err := store.RecordSample(context.Background(), value); err != nil {
		t.Fatal(err)
	}
	value.Tags["sequence"] = "changed"
	stored, err := store.Samples(context.Background(), "temperature", time.Time{}, time.Unix(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	if stored[0].Tags["sequence"] != "2" {
		t.Fatalf("stored sample follows caller map: %q", stored[0].Tags["sequence"])
	}
}

func TestSamplesReturnsIndependentTags(t *testing.T) {
	store := NewStore()
	store.samples = []metric.Sample{sampleForSnapshot(3)}
	first, err := store.Samples(context.Background(), "temperature", time.Time{}, time.Unix(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	first[0].Tags["sequence"] = "changed"
	second, err := store.Samples(context.Background(), "temperature", time.Time{}, time.Unix(10, 0))
	if err != nil {
		t.Fatal(err)
	}
	if second[0].Tags["sequence"] != "3" {
		t.Fatalf("Samples leaked mutable tags: %q", second[0].Tags["sequence"])
	}
}

func TestQuerySamplesReturnsIndependentTags(t *testing.T) {
	store := NewStore()
	store.samples = []metric.Sample{sampleForSnapshot(4)}
	first, err := store.QuerySamples(context.Background(), SampleQuery{MetricID: "temperature"})
	if err != nil {
		t.Fatal(err)
	}
	first[0].Tags["sequence"] = "changed"
	second, err := store.QuerySamples(context.Background(), SampleQuery{MetricID: "temperature"})
	if err != nil {
		t.Fatal(err)
	}
	if second[0].Tags["sequence"] != "4" {
		t.Fatalf("QuerySamples leaked mutable tags: %q", second[0].Tags["sequence"])
	}
}

func TestLatestReturnsIndependentTags(t *testing.T) {
	store := NewStore()
	store.samples = []metric.Sample{sampleForSnapshot(5)}
	first, ok := store.Latest(context.Background(), "temperature")
	if !ok {
		t.Fatal("latest sample missing")
	}
	first.Tags["sequence"] = "changed"
	second, ok := store.Latest(context.Background(), "temperature")
	if !ok {
		t.Fatal("latest sample missing on second read")
	}
	if second.Tags["sequence"] != "5" {
		t.Fatalf("Latest leaked mutable tags: %q", second.Tags["sequence"])
	}
}

func TestConcurrentRecordAndClearIsRaceFree(t *testing.T) {
	store := NewStore()
	for index := 0; index < 64; index++ {
		if err := store.RecordSample(context.Background(), sampleForSnapshot(index)); err != nil {
			t.Fatal(err)
		}
	}
	start := make(chan struct{})
	var wait sync.WaitGroup
	work := []func(){
		func() {
			for index := 64; index < 320; index++ {
				_ = store.RecordSample(context.Background(), sampleForSnapshot(index))
			}
		},
		func() {
			for index := 0; index < 256; index++ {
				_, _ = store.Samples(context.Background(), "temperature", time.Time{}, time.Unix(500, 0))
			}
		},
		func() {
			for index := 0; index < 256; index++ {
				_, _ = store.QuerySamples(context.Background(), SampleQuery{MetricID: "temperature"})
			}
		},
		func() {
			for index := 0; index < 256; index++ {
				_, _ = store.Latest(context.Background(), "temperature")
			}
		},
		func() {
			for index := 0; index < 256; index++ {
				store.ClearSamples(context.Background(), time.Unix(int64(index), 0))
			}
		},
	}
	for _, fn := range work {
		wait.Add(1)
		go func(run func()) {
			defer wait.Done()
			<-start
			run()
		}(fn)
	}
	close(start)
	wait.Wait()
}
