package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrEventNotFound = errors.New("event not found")
	ErrDateBusy      = errors.New("date busy")
)

type Event struct {
	ID          string    `db:"id"`
	Title       string    `db:"title"`
	Description string    `db:"description"`
	StartAt     time.Time `db:"start_at"`
	EndAt       time.Time `db:"end_at"`
	UserID      string    `db:"user_id"`
	Notified    bool      `db:"notified"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type EventFilter struct {
	From time.Time
	To   time.Time
}

type EventRepository interface {
	AddEvent(ctx context.Context, event Event) error
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventByID(ctx context.Context, id string) (Event, error)
	ListEvents(ctx context.Context) ([]Event, error)
	ListEventsWithFilter(ctx context.Context, filter EventFilter) ([]Event, error)
	GetEventsToNotify(ctx context.Context) ([]Event, error)
	MarkEventNotified(ctx context.Context, id string) error
	DeleteOldEvents(ctx context.Context, olderThan time.Time) (int64, error)
}
