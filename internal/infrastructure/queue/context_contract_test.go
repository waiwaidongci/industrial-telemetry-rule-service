package queue

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

var errQueueContextLost = errors.New("queue context was not forwarded")

func canceledQueueContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func receiveQueueResult(t *testing.T, result <-chan error) error {
	t.Helper()
	select {
	case err := <-result:
		return err
	case <-time.After(125 * time.Millisecond):
		return errQueueContextLost
	}
}

func TestMemoryPublishRejectsCanceledContext(t *testing.T) {
	queue := NewMemory()
	err := queue.Publish(canceledQueueContext(), Message{Key: "sample"})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Publish error = %v, want context.Canceled", err)
	}
	if queue.Len() != 0 {
		t.Fatalf("canceled publish queued %d messages", queue.Len())
	}
}

func TestMemoryConsumeForwardsHandlerContext(t *testing.T) {
	queue := NewMemory()
	if err := queue.Publish(context.Background(), Message{Key: "sample"}); err != nil {
		t.Fatal(err)
	}
	err := queue.Consume(canceledQueueContext(), func(ctx context.Context, _ Message) error {
		if err := ctx.Err(); err != nil {
			return err
		}
		return errQueueContextLost
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("handler context error = %v, want context.Canceled", err)
	}
}

func TestMemoryConsumeStopsWhileIdle(t *testing.T) {
	result := make(chan error, 1)
	go func() {
		result <- NewMemory().Consume(canceledQueueContext(), func(context.Context, Message) error { return nil })
	}()
	if err := receiveQueueResult(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("idle Consume error = %v, want context.Canceled", err)
	}
}

func TestRetryingConsumerSkipsCanceledHandler(t *testing.T) {
	var calls atomic.Int32
	consumer := &RetryingConsumer{Policy: RetryPolicy{Maximum: 2, InitialDelay: time.Second}}
	err := consumer.Handle(canceledQueueContext(), Message{Key: "sample"}, func(ctx context.Context, _ Message) error {
		calls.Add(1)
		if err := ctx.Err(); err != nil {
			return err
		}
		return errQueueContextLost
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Handle error = %v, want context.Canceled", err)
	}
	if calls.Load() != 0 {
		t.Fatalf("handler called %d times after cancellation", calls.Load())
	}
}

func TestRetryingConsumerStopsDuringBackoff(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	first := make(chan struct{})
	result := make(chan error, 1)
	consumer := &RetryingConsumer{Policy: RetryPolicy{Maximum: 3, InitialDelay: 500 * time.Millisecond}}
	go func() {
		result <- consumer.Handle(ctx, Message{Key: "sample"}, func(context.Context, Message) error {
			select {
			case <-first:
			default:
				close(first)
			}
			return errors.New("retry")
		})
	}()
	<-first
	cancel()
	if err := receiveQueueResult(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("backoff Handle error = %v, want context.Canceled", err)
	}
}

type blockingConsumer struct{ observed chan error }

func (c blockingConsumer) Consume(ctx context.Context, _ func(context.Context, Message) error) error {
	select {
	case <-ctx.Done():
		c.observed <- ctx.Err()
		return ctx.Err()
	case <-time.After(125 * time.Millisecond):
		c.observed <- errQueueContextLost
		return errQueueContextLost
	}
}

func TestWorkerPoolForwardsCancellationToConsumer(t *testing.T) {
	observed := make(chan error, 1)
	result := make(chan error, 1)
	pool := WorkerPool{Workers: 1, Consumer: blockingConsumer{observed: observed}, Handler: func(context.Context, Message) error { return nil }}
	go func() { result <- pool.Run(canceledQueueContext()) }()
	if err := receiveQueueResult(t, observed); !errors.Is(err, context.Canceled) {
		t.Fatalf("consumer context error = %v, want context.Canceled", err)
	}
	if err := receiveQueueResult(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("pool error = %v, want context.Canceled", err)
	}
}

func TestNopConsumerStopsOnCallerCancellation(t *testing.T) {
	result := make(chan error, 1)
	go func() {
		result <- (Nop{}).Consume(canceledQueueContext(), func(context.Context, Message) error { return nil })
	}()
	if err := receiveQueueResult(t, result); !errors.Is(err, context.Canceled) {
		t.Fatalf("Nop Consume error = %v, want context.Canceled", err)
	}
}
