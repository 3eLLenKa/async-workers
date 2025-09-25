package queue_test

import (
	"errors"
	"log/slog"
	"os"
	"testing"
	"workers/internal/queue"
)

func TestEnqueue(t *testing.T) {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	t.Helper()
	t.Parallel()

	t.Run("enqueue success", func(t *testing.T) {
		q := queue.New(log, 2) // очередь размером 2

		err := q.Enqueue("1", "payload1", 3)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		select {
		case task := <-q.Tasks():
			if task.Id != "1" || task.Payload != "payload1" || task.MaxRetries != 3 || task.Status != "queued" {
				t.Errorf("task fields mismatch: %+v", task)
			}
		default:
			t.Fatal("expected task in queue, but channel empty")
		}
	})

	t.Run("queue is full", func(t *testing.T) {
		q := queue.New(log, 1)

		if err := q.Enqueue("1", "payload1", 5); err != nil {
			t.Fatalf("unexpected error on first enqueue: %v", err)
		}

		err := q.Enqueue("2", "payload2", 5)

		if err == nil {
			t.Fatal("expected error when queue is full")
		}

		if !errors.Is(err, queue.ErrQueueIsFull) {
			t.Fatalf("expected ErrQueueIsFull, got: %v", err)
		}
	})
}
