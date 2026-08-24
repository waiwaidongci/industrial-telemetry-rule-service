package events

import (
	"testing"
	"time"

	"github.com/example/telemetry-rule-service/internal/domain/event"
)

func ownershipEvent(id, message string, at time.Time) event.Event {
	return event.Event{ID: id, RuleID: "rule", MetricID: "temperature", SourceID: "sensor", Message: message, Status: event.Open, FirstSeen: at, LastSeen: at, Count: 1, Labels: map[string]string{"site": id}}
}

func TestAggregateDoesNotMutateInputEvents(t *testing.T) {
	input := []event.Event{ownershipEvent("one", "alarm", time.Now())}
	result := Aggregate(input)
	result[0].Labels["site"] = "changed"
	if input[0].Labels["site"] != "one" {
		t.Fatalf("aggregate result shares labels with input: %#v", input[0].Labels)
	}
}

func TestApplyQueryDoesNotReuseInputStorage(t *testing.T) {
	now := time.Now()
	input := []event.Event{ownershipEvent("older", "alarm", now.Add(-time.Minute)), ownershipEvent("newer", "alarm", now)}
	result := ApplyQuery(input, Query{})
	if input[0].ID != "older" || input[1].ID != "newer" {
		t.Fatalf("query reordered caller slice: %#v", input)
	}
	result[0].Labels["site"] = "changed"
	if input[1].Labels["site"] != "newer" {
		t.Fatalf("query result shares labels with input: %#v", input[1].Labels)
	}
}

func TestSearchResultOwnsLabels(t *testing.T) {
	input := []event.Event{ownershipEvent("one", "pump alarm", time.Now())}
	for _, phrase := range []string{"pump", ""} {
		result := Search(input, phrase)
		result[0].Labels["site"] = "changed"
		if input[0].Labels["site"] != "one" {
			t.Fatalf("search %q shares labels with input: %#v", phrase, input[0].Labels)
		}
	}
}
