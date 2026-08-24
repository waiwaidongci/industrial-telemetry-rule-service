package event

import (
	"testing"
	"time"
)

func domainOwnershipEvent(id string, status Status, at time.Time) Event {
	return Event{ID: id, RuleID: "rule", MetricID: "temperature", SourceID: "sensor", Status: status, FirstSeen: at, LastSeen: at, Count: 1, Labels: map[string]string{"site": id}}
}

func TestFilterApplyOwnsResult(t *testing.T) {
	input := []Event{domainOwnershipEvent("one", Open, time.Now())}
	result := (Filter{Status: Open}).Apply(input)
	result[0].Labels["site"] = "changed"
	if input[0].Labels["site"] != "one" {
		t.Fatalf("filter result shares labels with input: %#v", input[0].Labels)
	}
}

func TestMergeCannotRewriteHistoricalLabels(t *testing.T) {
	now := time.Now()
	first := domainOwnershipEvent("first", Open, now)
	second := domainOwnershipEvent("second", Open, now.Add(time.Second))
	second.RuleID, second.MetricID, second.SourceID = first.RuleID, first.MetricID, first.SourceID
	second.Labels = map[string]string{"zone": "north"}
	merged := Merge(first, second)
	if _, exists := first.Labels["zone"]; exists {
		t.Fatalf("merge wrote into historical labels: %#v", first.Labels)
	}
	merged.Labels["site"] = "changed"
	if first.Labels["site"] != "first" {
		t.Fatalf("merged event still aliases historical labels: %#v", first.Labels)
	}
}

func TestBuildStatisticsDoesNotReorderEvents(t *testing.T) {
	now := time.Now()
	input := []Event{domainOwnershipEvent("newer", Open, now), domainOwnershipEvent("older", Resolved, now.Add(-time.Hour))}
	_ = BuildStatistics(input)
	if input[0].ID != "newer" || input[1].ID != "older" {
		t.Fatalf("statistics reordered caller events: %#v", input)
	}
}

func TestActiveOwnsResultLabels(t *testing.T) {
	now := time.Now()
	input := []Event{domainOwnershipEvent("closed", Resolved, now), domainOwnershipEvent("open", Open, now)}
	result := Active(input)
	if input[0].ID != "closed" {
		t.Fatalf("active filter overwrote caller slice: %#v", input)
	}
	result[0].Labels["site"] = "changed"
	if input[1].Labels["site"] != "open" {
		t.Fatalf("active result shares labels with input: %#v", input[1].Labels)
	}
}
