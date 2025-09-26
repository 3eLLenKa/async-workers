package main

import (
	"context"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"workers/internal/app"
	"workers/internal/config"
	"workers/internal/models"
	"workers/internal/queue"

	"log/slog"
)

const (
	statusRunning = "running"
	statusDone    = "done"
	statusFailed  = "failed"
)

func TaskWorker(ctx context.Context, log *slog.Logger, n int, tasks <-chan *models.Task) {
	log.Info("Worker running", "worker", n)

	for {
		select {
		case <-ctx.Done():
			log.Info("Worker is stopping", "worker", n)
			return
		case task, ok := <-tasks:
			if !ok {
				log.Info("Task channel closed", "worker", n)
				return
			}

			log.Info("Worker started the task", "worker", n, "id", task.Id)
			task.Status = statusRunning

			var attempt int

			for {
				attempt++ // считаем попытку

				workTime := rand.Intn(400) + 100 // симуляция работы 100–500 мс
				time.Sleep(time.Duration(workTime) * time.Millisecond)

				if rand.Float32() < 0.2 { // задача падает с вероятностью 20%
					if attempt > task.MaxRetries {
						task.Status = statusFailed

						log.Error("Task failed",
							"worker", n,
							"id", task.Id,
							"attempts", attempt,
						)

						break
					}

					// расчет backoff
					baseDelay := 100 * time.Millisecond
					maxDelay := 5 * time.Second

					backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt-1))) // экспонента от числа падений

					if backoff > maxDelay {
						backoff = maxDelay
					}

					// full jitter
					delay := time.Duration(rand.Int63n(int64(backoff)))

					log.Warn("Task failed, will retry",
						"worker", n,
						"id", task.Id,
						"attempt", attempt,
						"delay", delay,
					)

					time.Sleep(delay)

					continue
				}

				// успех
				task.Status = statusDone

				log.Info("Task completed successfully",
					"worker", n,
					"id", task.Id,
					"payload", task.Payload,
					"attempts", attempt,
				)

				break
			}
		}
	}
}

func main() {
	cfg := config.Load()

	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	q := queue.New(log, cfg.QueueSize) // создаем очередь

	application := app.New(":8080", q) // инициалиируем приложение (хендлеры, роутер, сервер)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		application.Server.Run() // поднимаем сервер
	}()

	wg := &sync.WaitGroup{}

	for i := 0; i < cfg.Workers; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			TaskWorker(ctx, log, n, q.Tasks())
		}(i + 1)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // ждем сигнал для завершения программы

	cancel()                     // отменяем контекст и останавливаем воркеры
	application.Server.Stop(ctx) // останавливаем сервер

	wg.Wait()
}
