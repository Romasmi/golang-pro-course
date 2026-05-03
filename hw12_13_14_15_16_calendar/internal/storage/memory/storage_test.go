package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	s := New()
	ctx := context.Background()

	event := domain.Event{
		ID:    "1",
		Title: "Test event",
	}

	t.Run("add event", func(t *testing.T) {
		err := s.AddEvent(ctx, event)
		require.NoError(t, err)

		err = s.AddEvent(ctx, event)
		require.ErrorIs(t, err, domain.ErrDateBusy)
	})

	t.Run("get event", func(t *testing.T) {
		got, err := s.GetEventByID(ctx, "1")
		require.NoError(t, err)
		require.Equal(t, event.ID, got.ID)
		require.Equal(t, event.Title, got.Title)
		require.False(t, got.CreatedAt.IsZero())
		require.False(t, got.UpdatedAt.IsZero())

		_, err = s.GetEventByID(ctx, "non-existent")
		require.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("update event", func(t *testing.T) {
		updated := event
		updated.Title = "Updated title"
		err := s.UpdateEvent(ctx, updated)
		require.NoError(t, err)

		got, err := s.GetEventByID(ctx, "1")
		require.NoError(t, err)
		require.Equal(t, "Updated title", got.Title)
		require.False(t, got.UpdatedAt.IsZero())

		err = s.UpdateEvent(ctx, domain.Event{ID: "non-existent"})
		require.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("list events", func(t *testing.T) {
		events, err := s.ListEvents(ctx)
		require.NoError(t, err)
		require.Len(t, events, 1)
	})

	t.Run("list events with filter", func(t *testing.T) {
		from := event.StartAt.Add(-time.Hour)
		to := event.StartAt.Add(time.Hour)
		filter := domain.EventFilter{From: from, To: to}
		events, err := s.ListEventsWithFilter(ctx, filter)
		require.NoError(t, err)
		require.Len(t, events, 1)

		filter = domain.EventFilter{From: to, To: to.Add(time.Hour)}
		events, err = s.ListEventsWithFilter(ctx, filter)
		require.NoError(t, err)
		require.Len(t, events, 0)
	})

	t.Run("delete event", func(t *testing.T) {
		err := s.DeleteEvent(ctx, "1")
		require.NoError(t, err)

		_, err = s.GetEventByID(ctx, "1")
		require.ErrorIs(t, err, domain.ErrEventNotFound)

		err = s.DeleteEvent(ctx, "1")
		require.ErrorIs(t, err, domain.ErrEventNotFound)
	})
}

func TestStorageConcurrency(t *testing.T) {
	s := New()
	ctx := context.Background()
	wg := sync.WaitGroup{}
	n := 100

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			id := fmt.Sprintf("event-%d", i)
			_ = s.AddEvent(ctx, domain.Event{ID: id})
		}(i)
	}
	wg.Wait()

	events, err := s.ListEvents(ctx)
	require.NoError(t, err)
	require.Len(t, events, n)
}
