package event

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Fingerprint(ruleID, sourceID, metricID string) string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%s", ruleID, sourceID, metricID)))
	return hex.EncodeToString(sum[:])
}
func GroupKey(e Event) string { return Fingerprint(e.RuleID, e.SourceID, e.MetricID) }
func Merge(a, b Event) Event {
	if a.ID == "" {
		return b
	}
	if b.LastSeen.After(a.LastSeen) {
		a.LastSeen = b.LastSeen
	}
	a.Count += b.Count
	if a.Message == "" {
		a.Message = b.Message
	}
	merged := make(map[string]string, len(a.Labels)+len(b.Labels))
	for key, value := range a.Labels {
		merged[key] = value
	}
	for key, value := range b.Labels {
		if _, exists := merged[key]; !exists {
			merged[key] = value
		}
	}
	if len(merged) > 0 {
		a.Labels = merged
	} else {
		a.Labels = nil
	}
	return a
}
