package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/event_service"
	memorystorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
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
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	config, err := NewConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logService := logger.New(config.Logger.Level)

	var st domain.EventRepository
	switch config.Storage.Type {
	case "memory":
		st = memorystorage.New()
	case "sql":
		s := sqlstorage.New()
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := s.Connect(ctx, config.Storage.DSN); err != nil {
			return fmt.Errorf("failed to connect to sql storage: %w", err)
		}
		if err := s.Migrate(ctx, config.Storage.MigrationsDir); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		st = s
	default:
		return fmt.Errorf("unknown storage type: %s", config.Storage.Type)
	}

	calendar := event_service.New(logService, st)

	ucs := usecases.NewUsecases(calendar)

	grpcSrv := grpcserver.NewServer(logService, ucs)

	grpcAddr := fmt.Sprintf("%s:%s", config.GRPC.Host, config.GRPC.Port)
	httpSrv := internalhttp.NewServer(logService, grpcAddr, config.HTTP.Host, config.HTTP.Port)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		if err := grpcSrv.Start(config.GRPC.Host, config.GRPC.Port); err != nil {
			logService.Error("failed to start grpc server: " + err.Error())
		}
	}()

	go func() {
		<-ctx.Done()

		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*3)
		defer stopCancel()

		if err := httpSrv.Stop(stopCtx); err != nil {
			logService.Error("failed to stop http server: " + err.Error())
		}
		grpcSrv.Stop()
	}()

	logService.Info("calendar is running...")

	if err := httpSrv.Start(ctx); err != nil {
		return fmt.Errorf("failed to start http server: %w", err)
	}

	return nil
}
