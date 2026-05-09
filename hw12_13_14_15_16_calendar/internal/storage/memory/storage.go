package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
)

type Storage struct {
	events map[string]domain.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]domain.Event),
	}
}

func (s *Storage) AddEvent(_ context.Context, event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[event.ID]; ok {
		return domain.ErrDateBusy
	}
	now := time.Now()
	event.CreatedAt = now
	event.UpdatedAt = now
	s.events[event.ID] = event
	return nil
}

func (s *Storage) UpdateEvent(_ context.Context, event domain.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	old, ok := s.events[event.ID]
	if !ok {
		return domain.ErrEventNotFound
	}
	event.CreatedAt = old.CreatedAt
	event.UpdatedAt = time.Now()
	s.events[event.ID] = event
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.events[id]; !ok {
		return domain.ErrEventNotFound
	}
	delete(s.events, id)
	return nil
}

func (s *Storage) GetEventByID(_ context.Context, id string) (domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	event, ok := s.events[id]
	if !ok {
		return domain.Event{}, domain.ErrEventNotFound
	}
	return event, nil
}

func (s *Storage) ListEvents(_ context.Context) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]domain.Event, 0, len(s.events))
	for _, e := range s.events {
		res = append(res, e)
	}
	return res, nil
}

func (s *Storage) ListEventsWithFilter(_ context.Context, filter domain.EventFilter) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]domain.Event, 0)
	for _, e := range s.events {
		if (e.StartAt.After(filter.From) || e.StartAt.Equal(filter.From)) && e.StartAt.Before(filter.To) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (s *Storage) GetEventsToNotify(_ context.Context) ([]domain.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	res := make([]domain.Event, 0)
	now := time.Now()
	for _, e := range s.events {
		if !e.Notified && (e.StartAt.Before(now) || e.StartAt.Equal(now)) {
			res = append(res, e)
		}
	}
	return res, nil
}

func (s *Storage) MarkEventNotified(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.events[id]
	if !ok {
		return domain.ErrEventNotFound
	}
	e.Notified = true
	s.events[id] = e
	return nil
}

func (s *Storage) DeleteOldEvents(_ context.Context, olderThan time.Time) (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	var count int64
	for id, e := range s.events {
		if e.StartAt.Before(olderThan) {
			delete(s.events, id)
			count++
		}
	}
	return count, nil
}
