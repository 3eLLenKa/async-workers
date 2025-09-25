package router

import (
	"net/http"
	"workers/internal/http/handlers"
)

type Router struct {
	handlers *handlers.Handlers
}

func New(handlers *handlers.Handlers) *Router {
	return &Router{
		handlers: handlers,
	}
}

func (r *Router) InitRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/queue/health", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.handlers.HealthCheckHandler(w, req)
	})

	mux.HandleFunc("/api/v1/queue/enqueue", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		r.handlers.EnqueueHandler(w, req)
	})

	return mux
}
