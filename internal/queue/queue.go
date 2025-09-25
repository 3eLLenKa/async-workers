package queue

import (
	"errors"
	"fmt"
	"log/slog"
	"workers/internal/models"
)

var ErrQueueIsFull = errors.New("queue is full")

type TaskQueue struct {
	log   *slog.Logger
	queue chan *models.Task
}

func New(log *slog.Logger, size int) *TaskQueue { // инициализация очереди
	return &TaskQueue{
		queue: make(chan *models.Task, size),
		log:   log,
	}
}

func (tq *TaskQueue) Enqueue(id, payload string, maxRetries int) error { // реализуем метод Enqueue
	log := tq.log.With(
		slog.String("id", id),
	)

	task := &models.Task{
		Id:         id,
		Payload:    payload,
		MaxRetries: maxRetries,
		Status:     "queued",
	}

	select {
	case tq.queue <- task:
		log.Info("task is stored")
		return nil
	default:
		log.Error(ErrQueueIsFull.Error())
		return fmt.Errorf("task %s: %w", id, ErrQueueIsFull)
	}
}

func (tq *TaskQueue) Tasks() <-chan *models.Task { // получение очереди задач
	return tq.queue
}
