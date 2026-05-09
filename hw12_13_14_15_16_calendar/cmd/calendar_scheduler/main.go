package main

import (
	"context"
	"flag"
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

	conf, err := scheduler.NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(conf.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := scheduler.New(conf, l)

	if err := app.Init(ctx); err != nil {
		l.Error("failed to init app: " + err.Error())
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		l.Error("failed to run app: " + err.Error())
		os.Exit(1)
	}
}
