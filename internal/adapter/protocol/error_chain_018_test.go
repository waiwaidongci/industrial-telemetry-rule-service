package protocol

import (
	"errors"
	"testing"

	"github.com/example/telemetry-rule-service/internal/domain/metric"
)

var errDecode018 = errors.New("decode sentinel 018")

type failingDecoder018 struct{}

func (failingDecoder018) Decode([]byte) (metric.Sample, error) { return metric.Sample{}, errDecode018 }

func TestRegistryPreservesDecodeCause018(t *testing.T) {
	r := NewRegistry()
	r.Register("application/test", failingDecoder018{})
	_, err := r.Decode("application/test", nil)
	if !errors.Is(err, errDecode018) { t.Fatalf("registry lost cause: %v", err) }
}

func TestRegistryPreservesUnsupportedCause018(t *testing.T) {
	_, err := NewRegistry().Decode("application/unknown", nil)
	if !errors.Is(err, ErrUnsupportedContentType) { t.Fatalf("unsupported type lost cause: %v", err) }
}
