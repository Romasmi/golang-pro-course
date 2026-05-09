package calendar

import (
	"context"
	"fmt"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	grpcserver "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/server/http"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/eventservice"
	memorystorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/usecases"
)

type App struct {
	config     Config
	logger     *logger.Logger
	service    *eventservice.CalendarService
	storage    domain.EventRepository
	grpcServer *grpcserver.Server
	httpServer *internalhttp.Server
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
		if err := s.Migrate(ctx, a.config.Storage.MigrationsDir); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}
		st = s
	default:
		return fmt.Errorf("unknown storage type: %s", a.config.Storage.Type)
	}
	a.service = eventservice.New(a.logger, st)
	a.storage = st

	ucs := usecases.NewUsecases(a.service)

	a.grpcServer = grpcserver.NewServer(a.logger, ucs)
	grpcAddr := fmt.Sprintf("%s:%s", a.config.GRPC.Host, a.config.GRPC.Port)
	a.httpServer = internalhttp.NewServer(a.logger, grpcAddr, a.config.HTTP.Host, a.config.HTTP.Port)

	return nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	go func() {
		if err := a.grpcServer.Start(a.config.GRPC.Host, a.config.GRPC.Port); err != nil {
			errCh <- fmt.Errorf("grpc server: %w", err)
		}
	}()

	go func() {
		if err := a.httpServer.Start(ctx); err != nil {
			errCh <- fmt.Errorf("http server: %w", err)
		}
	}()

	a.logger.Info("calendar is running...")

	select {
	case <-ctx.Done():
		a.logger.Info("Stopping servers...")
	case err := <-errCh:
		a.logger.Error("Server error: " + err.Error())
	}

	stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer stopCancel()

	if err := a.httpServer.Stop(stopCtx); err != nil {
		a.logger.Error("failed to stop http server: " + err.Error())
	}
	a.grpcServer.Stop()

	if closer, ok := a.storage.(interface{ Close(context.Context) error }); ok {
		if err := closer.Close(stopCtx); err != nil {
			a.logger.Error("failed to close storage: " + err.Error())
		}
	}

	return nil
}
