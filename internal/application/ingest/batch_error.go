package ingest

import (
	"fmt"
)

// BatchError is a JSON-serializable representation of a single batch item error.
// Unlike `error` (which encodes as `{}`), it preserves the failure reason
// so callers can inspect individual item failures after a round-trip.
type BatchError struct {
	Index int    `json:"index"`
	Error string `json:"error"`
}

func (e BatchError) String() string { return fmt.Sprintf("%d: %s", e.Index, e.Error) }
