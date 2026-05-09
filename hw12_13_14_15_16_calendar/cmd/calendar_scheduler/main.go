package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/app/scheduler"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "configs/scheduler_config.yaml", "path to config file")
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		log.Fatalf("%v", err)
	}
}

func run() error {
	conf, err := scheduler.NewConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	l := logger.New(conf.Logger.Level, "scheduler")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := scheduler.New(conf, l)

	if err := app.Init(ctx); err != nil {
		return fmt.Errorf("failed to init app: %w", err)
	}

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}
