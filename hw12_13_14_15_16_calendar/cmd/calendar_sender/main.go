package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/queue/rabbitmq"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/sender_config.yaml", "path to config file")
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

	// Connect to RabbitMQ
	rmq, err := rabbitmq.New(conf.RabbitMQ.URL, conf.RabbitMQ.Queue)
	if err != nil {
		l.Error("failed to connect to rabbitmq: " + err.Error())
		os.Exit(1)
	}
	defer rmq.Close()

	notifications, err := rmq.Receive(ctx)
	if err != nil {
		l.Error("failed to start receiving: " + err.Error())
		os.Exit(1)
	}

	l.Info("Sender started")

	for n := range notifications {
		msg := fmt.Sprintf("NOTIFICATION: EventID=%s, Title=%s, UserID=%s, StartAt=%s",
			n.EventID, n.Title, n.UserID, n.StartAt.Format("2006-01-02 15:04:05"))
		l.Info(msg)
	}

	l.Info("Sender stopping")
}
