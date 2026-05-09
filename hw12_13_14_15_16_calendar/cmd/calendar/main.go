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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	conf, err := calendar.NewConfig(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	l := logger.New(conf.Logger.Level)

	app := calendar.New(conf, l)

	if err := app.Init(ctx); err != nil {
		l.Error("failed to init app: " + err.Error())
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		l.Error("failed to run app: " + err.Error())
		os.Exit(1)
	}
}
