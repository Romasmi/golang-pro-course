package event_service

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
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	event := domain.Event{ID: "1", Title: "Test"}
	repo.On("AddEvent", mock.Anything, event).Return(nil)

	res, err := s.CreateEvent(context.Background(), event)
	assert.NoError(t, err)
	assert.Equal(t, event, res)
	repo.AssertExpectations(t)
}

func TestCalendarService_GetEvent(t *testing.T) {
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

func TestCalendarService_ListEventsByFilter(t *testing.T) {
	repo := new(MockRepository)
	logger := new(MockLogger)
	s := New(logger, repo)

	from := time.Now()
	to := from.Add(time.Hour)
	filter := domain.EventFilter{From: from, To: to}
	events := []domain.Event{{ID: "1", Title: "Test"}}

	repo.On("ListEventsWithFilter", mock.Anything, filter).Return(events, nil)

	res, err := s.ListEventsByFilter(context.Background(), from, to)
	assert.NoError(t, err)
	assert.Equal(t, events, res)
	repo.AssertExpectations(t)
}
