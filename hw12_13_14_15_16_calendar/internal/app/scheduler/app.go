package scheduler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
	memorystorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/queue"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/queue/rabbitmq"
)

type App struct {
	config     Config
	logger     *logger.Logger
	service    *eventservice.CalendarService
	storage    domain.EventRepository
	rmq        *rabbitmq.RabbitMQ
	httpServer *http.Server
}

func New(conf Config, l *logger.Logger) *App {
	return &App{
		config: conf,
		logger: l,
	}
}

func (a *App) Init(ctx context.Context) error {
	var st domain.EventRepository
	switch a.config.Storage.Type {
	case "memory":
		st = memorystorage.New()
	case "sql":
		s := sqlstorage.New()
		if err := s.Connect(ctx, a.config.Storage.DSN); err != nil {
			return fmt.Errorf("failed to connect to sql storage: %w", err)
		}
		st = s
	default:
		return fmt.Errorf("unknown storage type: %s", a.config.Storage.Type)
	}
	a.service = eventservice.New(a.logger, st)
	a.storage = st

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
	a.logger.Info("Scheduler started")

	go func() {
		if a.config.HTTP.Port != "" {
			a.logger.Info("Starting health server on " + a.httpServer.Addr)
			if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				a.logger.Error("health server failed: " + err.Error())
			}
		}
	}()

	scanTicker := time.NewTicker(a.config.Scheduler.ScanInterval)
	defer scanTicker.Stop()

	cleanTicker := time.NewTicker(a.config.Scheduler.CleanInterval)
	defer cleanTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			a.logger.Info("Scheduler stopping")
			stopCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			if a.httpServer != nil {
				if err := a.httpServer.Shutdown(stopCtx); err != nil {
					a.logger.Error("failed to stop health server: " + err.Error())
				}
			}
			a.rmq.Close()
			if closer, ok := a.storage.(interface{ Close(context.Context) error }); ok {
				if err := closer.Close(stopCtx); err != nil {
					a.logger.Error("failed to close storage: " + err.Error())
				}
			}
			return nil
		case <-scanTicker.C:
			a.scanEvents(ctx)
		case <-cleanTicker.C:
			a.cleanOldEvents(ctx)
		}
	}
}

func (a *App) scanEvents(ctx context.Context) {
	events, err := a.service.GetEventsToNotify(ctx)
	if err != nil {
		a.logger.Error("failed to get events to notify: " + err.Error())
		return
	}

	for _, e := range events {
		select {
		case <-ctx.Done():
			return
		default:
		}
		n := queue.Notification{
			EventID: e.ID,
			Title:   e.Title,
			StartAt: e.StartAt,
			UserID:  e.UserID,
		}

		if err := a.rmq.Publish(ctx, n); err != nil {
			a.logger.Error("failed to publish notification for event " + e.ID + ": " + err.Error())
			continue
		}

		if err := a.service.MarkEventNotified(ctx, e.ID); err != nil {
			a.logger.Error("failed to mark event as notified " + e.ID + ": " + err.Error())
		} else {
			a.logger.Info("published notification and marked as notified: " + e.ID)
		}
	}
}

func (a *App) cleanOldEvents(ctx context.Context) {
	olderThan := time.Now().AddDate(-1, 0, 0)
	count, err := a.service.DeleteOldEvents(ctx, olderThan)
	if err != nil {
		a.logger.Error("failed to clean old events: " + err.Error())
		return
	}
	if count > 0 {
		a.logger.Info(fmt.Sprintf("cleaned %d old events", count))
	}
}
