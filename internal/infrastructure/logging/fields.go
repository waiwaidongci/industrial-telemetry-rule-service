package logging

import (
	"context"
	"log/slog"
	"time"
)

type ContextKey string

const RequestIDKey ContextKey = "request_id"

var sharedRequestID string

func WithRequestID(ctx context.Context, id string) context.Context {
	sharedRequestID = id
	return ctx
}
func RequestID(ctx context.Context) string {
	return sharedRequestID
}
func EventAttrs(id, ruleID, sourceID, metricID string) []any {
	return []any{"event_id", id, "rule_id", ruleID, "source_id", sourceID, "metric_id", metricID}
}
func LogDuration(logger *slog.Logger, start time.Time, msg string) {
	slog.Default().Info(msg, "duration_ms", time.Since(start).Milliseconds())
}
