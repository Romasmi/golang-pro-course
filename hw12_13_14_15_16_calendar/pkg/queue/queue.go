package queue

import (
	"context"
	"time"
)

type Notification struct {
	EventID string    `json:"event_id"`
	Title   string    `json:"title"`
	StartAt time.Time `json:"start_at"`
	UserID  string    `json:"user_id"`
}

type Sender interface {
	Publish(ctx context.Context, n Notification) error
	Close() error
}

type Receiver interface {
	Receive(ctx context.Context) (<-chan Notification, error)
	Close() error
}
