package eventservice

import (
	"context"
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) AddEvent(ctx context.Context, event domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockRepository) UpdateEvent(ctx context.Context, event domain.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockRepository) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetEventByID(ctx context.Context, id string) (domain.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(domain.Event), args.Error(1)
}

func (m *MockRepository) ListEvents(ctx context.Context) ([]domain.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockRepository) ListEventsWithFilter(ctx context.Context, filter domain.EventFilter) ([]domain.Event, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockRepository) GetEventsToNotify(ctx context.Context) ([]domain.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]domain.Event), args.Error(1)
}

func (m *MockRepository) MarkEventNotified(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) DeleteOldEvents(ctx context.Context, olderThan time.Time) (int64, error) {
	args := m.Called(ctx, olderThan)
	return args.Get(0).(int64), args.Error(1)
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string) {
	m.Called(msg)
}

func (m *MockLogger) Error(msg string) {
	m.Called(msg)
}

func TestCalendarService_CreateEvent(t *testing.T) {
	t.Parallel()

	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		event := domain.Event{Title: "Test"}
		repo.On("AddEvent", mock.Anything, mock.MatchedBy(func(e domain.Event) bool {
			return e.Title == "Test" && e.ID != ""
		})).Return(nil).Once()

		res, err := s.CreateEvent(context.Background(), event)
		assert.NoError(t, err)
		assert.Equal(t, "Test", res.Title)
		assert.NotEmpty(t, res.ID)
	})

	repo.AssertExpectations(t)
}

func TestCalendarService_UpdateEvent(t *testing.T) {
	t.Parallel()

	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	event := domain.Event{ID: "1", Title: "Updated"}
	repo.On("UpdateEvent", mock.Anything, event).Return(nil)

	res, err := s.UpdateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, event, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_GetEvent(t *testing.T) {
	t.Parallel()

	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	event := domain.Event{ID: "1", Title: "Test"}
	repo.On("GetEventByID", mock.Anything, "1").Return(event, nil)

	res, err := s.GetEvent(context.Background(), "1")
	assert.NoError(t, err)
	assert.Equal(t, event, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_ListEvents(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	events := []domain.Event{{ID: "1", Title: "Test"}}
	repo.On("ListEvents", mock.Anything).Return(events, nil)

	res, err := s.ListEvents(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, events, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_ListEventsByFilter(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	from := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2026, 5, 2, 0, 0, 0, 0, time.UTC)
	filter := domain.EventFilter{From: from, To: to}
	events := []domain.Event{{ID: "1", Title: "Test", StartAt: from.Add(time.Hour)}}

	repo.On("ListEventsWithFilter", mock.Anything, filter).Return(events, nil)

	res, err := s.ListEventsByFilter(context.Background(), from, to)
	assert.NoError(t, err)
	assert.Equal(t, events, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_DeleteEvent(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	repo.On("DeleteEvent", mock.Anything, "1").Return(nil)

	err := s.DeleteEvent(context.Background(), "1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCalendarService_GetEventsToNotify(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	events := []domain.Event{{ID: "1", Title: "Test"}}
	repo.On("GetEventsToNotify", mock.Anything).Return(events, nil)

	res, err := s.GetEventsToNotify(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, events, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_MarkEventNotified(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	repo.On("MarkEventNotified", mock.Anything, "1").Return(nil)

	err := s.MarkEventNotified(context.Background(), "1")
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestCalendarService_DeleteOldEvents(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	now := time.Now()
	repo.On("DeleteOldEvents", mock.Anything, now).Return(int64(5), nil)

	count, err := s.DeleteOldEvents(context.Background(), now)
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	repo.AssertExpectations(t)
}
