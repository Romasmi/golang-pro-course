package sender

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/queue/rabbitmq"
)

type App struct {
	config Config
	logger *logger.Logger
	rmq    *rabbitmq.RabbitMQ
}

func New(conf Config, l *logger.Logger) *App {
	return &App{
		config: conf,
		logger: l,
	}
}

func (a *App) Init(ctx context.Context) error {
	rmq, err := rabbitmq.New(a.config.RabbitMQ.URL, a.config.RabbitMQ.Queue, a.logger)
	if err != nil {
		return fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	a.rmq = rmq

	return nil
}

func (a *App) Run(ctx context.Context) error {
	a.logger.Info("Sender started")

	notifications, err := a.rmq.Receive(ctx)
	if err != nil {
		return fmt.Errorf("failed to start receiving: %w", err)
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Sender stopping")
			a.rmq.Close()
			return nil
		case n, ok := <-notifications:
			if !ok {
				a.logger.Info("Notifications channel closed")
				return nil
			}
			timer := time.NewTimer(time.Duration(
				r.Intn(1900)+100,
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
