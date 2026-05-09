package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	conf, err := calendar.NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	l := logger.New(conf.Logger.Level, "calendar")

	app := calendar.New(conf, l)

	if err := app.Init(ctx); err != nil {
		return fmt.Errorf("failed to init app: %w", err)
	}

	if err := app.Run(ctx); err != nil {
		return fmt.Errorf("failed to run app: %w", err)
	}

	return nil
}
