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
		s := New()
		t1 := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
		t2 := time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)
		t3 := time.Date(2026, 5, 3, 10, 0, 0, 0, time.UTC)

		_ = s.AddEvent(ctx, domain.Event{ID: "1", Title: "E1", StartAt: t1})
		_ = s.AddEvent(ctx, domain.Event{ID: "2", Title: "E2", StartAt: t2})
		_ = s.AddEvent(ctx, domain.Event{ID: "3", Title: "E3", StartAt: t3})

		tests := []struct {
			name string
			from time.Time
			to   time.Time
			want int
		}{
			{
				name: "all events",
				from: t1.Add(-time.Hour),
				to:   t3.Add(time.Hour),
				want: 3,
			},
			{
				name: "first two",
				from: t1.Add(-time.Hour),
				to:   t2.Add(time.Hour),
				want: 2,
			},
			{
				name: "only second",
				from: t2.Add(-time.Hour),
				to:   t2.Add(time.Hour),
				want: 1,
			},
			{
				name: "none (before)",
				from: t1.Add(-2 * time.Hour),
				to:   t1.Add(-time.Hour),
				want: 0,
			},
			{
				name: "none (after)",
				from: t3.Add(time.Hour),
				to:   t3.Add(2 * time.Hour),
				want: 0,
			},
			{
				name: "boundary check [from, to)",
				from: t1,
				to:   t2,
				want: 1, // only E1
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				filter := domain.EventFilter{From: tt.from, To: tt.to}
				events, err := s.ListEventsWithFilter(ctx, filter)
				require.NoError(t, err)
				require.Len(t, events, tt.want)
			})
		}
	})

	t.Run("delete event", func(t *testing.T) {
		err := s.DeleteEvent(ctx, "1")
		require.NoError(t, err)

		_, err = s.GetEventByID(ctx, "1")
		require.ErrorIs(t, err, domain.ErrEventNotFound)

		err = s.DeleteEvent(ctx, "1")
		require.ErrorIs(t, err, domain.ErrEventNotFound)
	})

	t.Run("notify and clean", func(t *testing.T) {
		s := New()
		now := time.Now()

		// Should notify (StartAt - 10s <= now)
		_ = s.AddEvent(ctx, domain.Event{ID: "n1", StartAt: now.Add(5 * time.Second), RemindBefore: 10})
		// Should not notify (StartAt - 10s > now)
		_ = s.AddEvent(ctx, domain.Event{ID: "n2", StartAt: now.Add(20 * time.Second), RemindBefore: 10})
		// Already notified
		_ = s.AddEvent(ctx, domain.Event{ID: "n3", StartAt: now.Add(5 * time.Second), RemindBefore: 10, Notified: true})
		// No reminder
		_ = s.AddEvent(ctx, domain.Event{ID: "n4", StartAt: now.Add(5 * time.Second), RemindBefore: 0})

		toNotify, err := s.GetEventsToNotify(ctx)
		require.NoError(t, err)
		require.Len(t, toNotify, 1)
		require.Equal(t, "n1", toNotify[0].ID)

		err = s.MarkEventNotified(ctx, "n1")
		require.NoError(t, err)

		toNotify, _ = s.GetEventsToNotify(ctx)
		require.Len(t, toNotify, 0)

		// Clean old
		_ = s.AddEvent(ctx, domain.Event{ID: "old", StartAt: now.AddDate(-2, 0, 0)})
		count, err := s.DeleteOldEvents(ctx, now.AddDate(-1, 0, 0))
		require.NoError(t, err)
		require.Equal(t, int64(1), count)

		_, err = s.GetEventByID(ctx, "old")
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
