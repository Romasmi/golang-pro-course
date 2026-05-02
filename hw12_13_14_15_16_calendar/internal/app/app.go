package app

import (
	"context"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
)

type App struct {
	logger  Logger
	storage domain.EventRepository
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(logger Logger, storage domain.EventRepository) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}
