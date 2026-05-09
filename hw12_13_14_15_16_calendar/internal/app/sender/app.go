package sender

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/queue/rabbitmq"
)

type App struct {
	config     Config
	logger     *logger.Logger
	rmq        *rabbitmq.RabbitMQ
	httpServer *http.Server
}

func New(conf Config, l *logger.Logger) *App {
	return &App{
		config: conf,
		logger: l,
	}
}

func (a *App) Init(_ context.Context) error {
	rmq, err := rabbitmq.New(a.config.RabbitMQ.URL, a.config.RabbitMQ.Queue, a.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	a.rmq = rmq

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	a.httpServer = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", a.config.HTTP.Host, a.config.HTTP.Port),
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
	}

	return nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Sender started")

	go func() {
		if a.config.HTTP.Port != "" {
			a.logger.Info("Starting health server on " + a.httpServer.Addr)
			if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				a.logger.Error("health server failed: " + err.Error())
			}
		}
	}()

	notifications, err := a.rmq.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to start receiving: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Sender stopping")
			stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if a.httpServer != nil {
				if err := a.httpServer.Shutdown(stopCtx); err != nil {
					a.logger.Error("failed to stop health server: " + err.Error())
				}
			}
			a.rmq.Close()
			return nil
		case n, ok := <-notifications:
			if !ok {
				a.logger.Info("Notifications channel closed")
				return nil
			}
			timer := time.NewTimer(time.Duration(
				rand.IntN(1900)+100, //nolint:gosec
			) * time.Millisecond)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil
			case <-timer.C:
			}
			a.logger.Info(fmt.Sprintf("SENDER: Sending notification for event %s (user %s): %s",
				n.EventID, n.UserID, n.Title))
		}
	}
}
