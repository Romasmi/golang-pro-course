package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockEventServiceServer struct {
	api.UnimplementedEventServiceServer
	mock.Mock
}

func (m *mockEventServiceServer) CreateEvent(
	ctx context.Context,
	req *api.CreateEventRequest,
) (*api.CreateEventResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.CreateEventResponse), args.Error(1)
}

func (m *mockEventServiceServer) UpdateEvent(
	ctx context.Context,
	req *api.UpdateEventRequest,
) (*api.UpdateEventResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.UpdateEventResponse), args.Error(1)
}

func (m *mockEventServiceServer) GetEvent(
	ctx context.Context,
	req *api.GetEventRequest,
) (*api.GetEventResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.GetEventResponse), args.Error(1)
}

func (m *mockEventServiceServer) ListEvents(
	ctx context.Context,
	req *api.ListEventsRequest,
) (*api.ListEventsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.ListEventsResponse), args.Error(1)
}

func (m *mockEventServiceServer) ListEventsByInterval(
	ctx context.Context,
	req *api.ListEventsByIntervalRequest,
) (*api.ListEventsResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.ListEventsResponse), args.Error(1)
}

func (m *mockEventServiceServer) Hello(ctx context.Context, req *api.HelloRequest) (*api.HelloResponse, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*api.HelloResponse), args.Error(1)
}

func TestAPI(t *testing.T) {
	ctx := context.Background()
	gwmux := runtime.NewServeMux()
	mockServer := &mockEventServiceServer{}

	err := api.RegisterEventServiceHandlerServer(ctx, gwmux, mockServer)
	require.NoError(t, err)

	mux := http.NewServeMux()
	mux.Handle("/", gwmux)

	t.Run("CreateEvent", func(t *testing.T) {
		event := &api.Event{Title: "New Event", UserId: "user1"}
		expectedResponse := &api.CreateEventResponse{Event: &api.Event{Id: "1", Title: "New Event", UserId: "user1"}}
		mockServer.On("CreateEvent", mock.Anything, mock.MatchedBy(func(req *api.CreateEventRequest) bool {
			return req.Event.Title == event.Title && req.Event.UserId == event.UserId
		})).Return(expectedResponse, nil).Once()

		body, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.CreateEventResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, "1", resp.Event.Id)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		event := &api.Event{Id: "1", Title: "Updated Event"}
		expectedResponse := &api.UpdateEventResponse{Event: event}
		mockServer.On("UpdateEvent", mock.Anything, mock.MatchedBy(func(req *api.UpdateEventRequest) bool {
			return req.Event.Id == event.Id && req.Event.Title == event.Title
		})).Return(expectedResponse, nil).Once()

		body, _ := json.Marshal(event)
		req := httptest.NewRequest(http.MethodPut, "/events/1", bytes.NewReader(body))
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.UpdateEventResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, "Updated Event", resp.Event.Title)
	})

	t.Run("GetEvent", func(t *testing.T) {
		eventID := "test-id"
		expectedResponse := &api.GetEventResponse{
			Event: &api.Event{Id: eventID, Title: "Test Event"},
		}
		mockServer.On("GetEvent", mock.Anything, &api.GetEventRequest{Id: eventID}).Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/events/"+eventID, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.GetEventResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, eventID, resp.Event.Id)
	})

	t.Run("ListEvents", func(t *testing.T) {
		expectedResponse := &api.ListEventsResponse{
			Events: []*api.Event{{Id: "1", Title: "E1"}},
		}
		mockServer.On("ListEvents", mock.Anything, &api.ListEventsRequest{}).Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/events", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.ListEventsResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Len(t, resp.Events, 1)
	})

	t.Run("ListEventsByInterval", func(t *testing.T) {
		expectedResponse := &api.ListEventsResponse{
			Events: []*api.Event{{Id: "1", Title: "E1"}},
		}
		// The interval in path is 'DAY', 'WEEK', etc. gRPC gateway maps it to enum
		mockServer.On("ListEventsByInterval", mock.Anything, mock.MatchedBy(func(req *api.ListEventsByIntervalRequest) bool {
			return req.Interval == api.Interval_DAY
		})).Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/events/interval/DAY", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.ListEventsResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Len(t, resp.Events, 1)
	})

	t.Run("Hello", func(t *testing.T) {
		expectedResponse := &api.HelloResponse{Message: "Hi"}
		mockServer.On("Hello", mock.Anything, mock.Anything).Return(expectedResponse, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		require.Equal(t, http.StatusOK, rr.Code)
		var resp api.HelloResponse
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		require.Equal(t, "Hi", resp.Message)
	})
}
