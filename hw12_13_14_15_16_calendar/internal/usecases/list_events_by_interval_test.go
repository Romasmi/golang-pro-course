package usecases

import (
	"context"
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/services/event_service"
	memorystorage "github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockLogger struct{}

func (m *mockLogger) Info(msg string)  {}
func (m *mockLogger) Error(msg string) {}

func TestListEventsByIntervalUsecase_Do(t *testing.T) {
	st := memorystorage.New()
	svc := event_service.New(&mockLogger{}, st)
	uc := &ListEventsByIntervalUsecase{Service: svc}
	ctx := context.Background()

	// Events at different dates
	e1 := domain.Event{ID: "1", Title: "Day Event", StartAt: time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)}
	e2 := domain.Event{ID: "2", Title: "Next Day Event", StartAt: time.Date(2026, 5, 2, 10, 0, 0, 0, time.UTC)}
	e3 := domain.Event{ID: "3", Title: "Next Month Event", StartAt: time.Date(2026, 6, 1, 10, 0, 0, 0, time.UTC)}
	e4 := domain.Event{ID: "4", Title: "Next Year Event", StartAt: time.Date(2027, 1, 1, 10, 0, 0, 0, time.UTC)}

	_ = st.AddEvent(ctx, e1)
	_ = st.AddEvent(ctx, e2)
	_ = st.AddEvent(ctx, e3)
	_ = st.AddEvent(ctx, e4)

	tests := []struct {
		name    string
		req     ListEventsByIntervalRequest
		wantLen int
		wantIDs []string
	}{
		{
			name: "day interval - matches e1",
			req: ListEventsByIntervalRequest{
				Date:     time.Date(2026, 5, 1, 15, 0, 0, 0, time.UTC),
				Interval: Day,
			},
			wantLen: 1,
			wantIDs: []string{"1"},
		},
		{
			name: "week interval - matches e1, e2",
			req: ListEventsByIntervalRequest{
				Date:     time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
				Interval: Week,
			},
			wantLen: 2,
			wantIDs: []string{"1", "2"},
		},
		{
			name: "month interval - matches e1, e2",
			req: ListEventsByIntervalRequest{
				Date:     time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC),
				Interval: Month,
			},
			wantLen: 2,
			wantIDs: []string{"1", "2"},
		},
		{
			name: "year interval - matches e1, e2, e3",
			req: ListEventsByIntervalRequest{
				Date:     time.Date(2026, 12, 31, 0, 0, 0, 0, time.UTC),
				Interval: Year,
			},
			wantLen: 3,
			wantIDs: []string{"1", "2", "3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := uc.Do(ctx, tt.req)
			require.NoError(t, err)

			events := res.([]domain.Event)
			assert.Len(t, events, tt.wantLen)

			ids := make([]string, 0, len(events))
			for _, e := range events {
				ids = append(ids, e.ID)
			}
			for _, id := range tt.wantIDs {
				assert.Contains(t, ids, id)
			}
		})
	}
}
