package app

import (
	"workers/internal/http/handlers"
	"workers/internal/http/router"
	"workers/internal/http/server"
	"workers/internal/queue"
)

type App struct {
	Server *server.Server
}

func New(
	serverAddr string,
	q *queue.TaskQueue,
) *App {
	handlers := handlers.New(q)
	router := router.New(handlers).InitRouter()
	app := server.New(serverAddr, router)

	return &App{
		Server: app,
	}
}
