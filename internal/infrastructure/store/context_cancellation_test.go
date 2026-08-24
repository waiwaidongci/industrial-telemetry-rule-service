package store

import (
	"context"
	"errors"
	"testing"
)

func canceledContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	return ctx
}

func TestHealthPingHonorsCanceledContext(t *testing.T) {
	err := NewHealth().Ping(canceledContext())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Ping() error = %v, want context.Canceled", err)
	}
}

func TestPostgreSQLPingHonorsCanceledContext(t *testing.T) {
	err := NewPostgreSQLStore("postgres://telemetry").Ping(canceledContext())
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Ping() error = %v, want context.Canceled", err)
	}
}

func TestMigrationUpHonorsCanceledContext(t *testing.T) {
	runner := NewMigrationRunner()
	err := runner.Up(canceledContext(), []Migration{{Version: 1, Name: "initial"}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Up() error = %v, want context.Canceled", err)
	}
	if len(runner.Applied) != 0 {
		t.Fatalf("Up() applied %v after cancellation", runner.Applied)
	}
}

func TestMigrationDownHonorsCanceledContext(t *testing.T) {
	runner := NewMigrationRunner()
	runner.Applied[1] = true
	err := runner.Down(canceledContext(), []Migration{{Version: 1, Name: "initial"}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Down() error = %v, want context.Canceled", err)
	}
	if !runner.Applied[1] {
		t.Fatal("Down() removed migration after cancellation")
	}
}
