package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"workers/internal/queue"
)

type Worker interface {
	Enqueue(id, payload string, maxRetries int) error
}

type Handlers struct {
	worker Worker
}

func New(worker Worker) *Handlers {
	return &Handlers{
		worker: worker,
	}
}

func (h *Handlers) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 OK"))
}

func (h *Handlers) EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id         string `json:"id"`
		Payload    string `json:"payload"`
		MaxRetries int    `json:"max_retries"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.worker.Enqueue(req.Id, req.Payload, req.MaxRetries); err != nil {
		if errors.Is(err, queue.ErrQueueIsFull) {
			http.Error(w, err.Error(), http.StatusTooManyRequests)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusCreated)
}
