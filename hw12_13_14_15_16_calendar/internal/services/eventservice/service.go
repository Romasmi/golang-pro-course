package eventservice

import (
	"context"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/google/uuid"
)

type CalendarService struct {
	logger  Logger
	storage domain.EventRepository
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

func New(logger Logger, storage domain.EventRepository) *CalendarService {
	return &CalendarService{
		logger:  logger,
		storage: storage,
	}
}

func (s *CalendarService) CreateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	event.ID = uuid.New().String()
	err := s.storage.AddEvent(ctx, event)
	return event, err
}

func (s *CalendarService) UpdateEvent(ctx context.Context, event domain.Event) (domain.Event, error) {
	err := s.storage.UpdateEvent(ctx, event)
	return event, err
}

func (s *CalendarService) GetEvent(ctx context.Context, id string) (domain.Event, error) {
	return s.storage.GetEventByID(ctx, id)
}

func (s *CalendarService) ListEvents(ctx context.Context) ([]domain.Event, error) {
	return s.storage.ListEvents(ctx)
}

func (s *CalendarService) ListEventsByFilter(ctx context.Context, from, to time.Time) ([]domain.Event, error) {
	return s.storage.ListEventsWithFilter(ctx, domain.EventFilter{From: from, To: to})
}
