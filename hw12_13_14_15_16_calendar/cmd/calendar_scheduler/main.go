package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/queue"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	sqlstorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/sql"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/scheduler_config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	conf, err := NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(conf.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Connect to Storage
	storage := sqlstorage.New()
	if err := storage.Connect(ctx, conf.Storage.DSN); err != nil {
		l.Error("failed to connect to storage: " + err.Error())
		os.Exit(1)
	}
	defer storage.Close(ctx)

	// Connect to RabbitMQ
	rmq, err := rabbitmq.New(conf.RabbitMQ.URL, conf.RabbitMQ.Queue)
	if err != nil {
		l.Error("failed to connect to rabbitmq: " + err.Error())
		os.Exit(1)
	}
	defer rmq.Close()

	l.Info("Scheduler started")

	scanTicker := time.NewTicker(conf.Scheduler.ScanInterval)
	defer scanTicker.Stop()

	cleanTicker := time.NewTicker(conf.Scheduler.CleanInterval)
	defer cleanTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			l.Info("Scheduler stopping")
			return
		case <-scanTicker.C:
			scanEvents(ctx, storage, rmq, l)
		case <-cleanTicker.C:
			cleanOldEvents(ctx, storage, l)
		}
	}
}

func scanEvents(ctx context.Context, storage *sqlstorage.Storage, sender queue.Sender, l *logger.Logger) {
	events, err := storage.GetEventsToNotify(ctx)
	if err != nil {
		l.Error("failed to get events to notify: " + err.Error())
		return
	}

	for _, e := range events {
		n := queue.Notification{
			EventID: e.ID,
			Title:   e.Title,
			StartAt: e.StartAt,
			UserID:  e.UserID,
		}

		if err := sender.Publish(ctx, n); err != nil {
			l.Error("failed to publish notification for event " + e.ID + ": " + err.Error())
			continue
		}

		if err := storage.MarkEventNotified(ctx, e.ID); err != nil {
			l.Error("failed to mark event as notified " + e.ID + ": " + err.Error())
		} else {
			l.Info("published notification and marked as notified: " + e.ID)
		}
	}
}

func cleanOldEvents(ctx context.Context, storage *sqlstorage.Storage, l *logger.Logger) {
	olderThan := time.Now().AddDate(-1, 0, 0)
	count, err := storage.DeleteOldEvents(ctx, olderThan)
	if err != nil {
		l.Error("failed to clean old events: " + err.Error())
		return
	}
	if count > 0 {
		l.Info(fmt.Sprintf("cleaned %d old events", count))
	}
}
