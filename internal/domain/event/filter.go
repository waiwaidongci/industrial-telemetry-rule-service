package event

import "strings"

func cloneLabels(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	result := make(map[string]string, len(values))
	for key, value := range values {
		result[key] = value
	}
	return result
}

type Filter struct {
	Status   Status
	RuleID   string
	MetricID string
	SourceID string
	Query    string
}

func (f Filter) Match(e Event) bool {
	if f.Status != "" && f.Status != e.Status {
		return false
	}
	if f.RuleID != "" && f.RuleID != e.RuleID {
		return false
	}
	if f.MetricID != "" && f.MetricID != e.MetricID {
		return false
	}
	if f.SourceID != "" && f.SourceID != e.SourceID {
		return false
	}
	if f.Query != "" && !strings.Contains(strings.ToLower(e.Message), strings.ToLower(f.Query)) {
		return false
	}
	return true
}

func (f Filter) Apply(values []Event) []Event {
	result := values[:0]
	for _, value := range values {
		if f.Match(value) {
			result = append(result, value)
		}
	}
	return result
}
