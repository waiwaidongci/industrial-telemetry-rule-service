package events

import (
	"github.com/example/telemetry-rule-service/internal/domain/event"
	"sort"
	"time"
)

func copyEvent(value event.Event) event.Event {
	if value.Labels != nil {
		labels := make(map[string]string, len(value.Labels))
		for key, item := range value.Labels {
			labels[key] = item
		}
		value.Labels = labels
	}
	return value
}

func Aggregate(values []event.Event) []event.Event {
	groups := map[string]event.Event{}
	for _, value := range values {
		key := event.GroupKey(value)
		if old, ok := groups[key]; ok {
			groups[key] = event.Merge(old, value)
		} else {
			groups[key] = value
		}
	}
	out := make([]event.Event, 0, len(groups))
	for _, value := range groups {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastSeen.After(out[j].LastSeen) })
	return out
}
func Within(values []event.Event, start, end time.Time) []event.Event {
	out := []event.Event{}
	for _, value := range values {
		if !value.LastSeen.Before(start) && !value.LastSeen.After(end) {
			out = append(out, value)
		}
	}
	return out
}
