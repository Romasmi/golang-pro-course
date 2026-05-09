package integration

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/pkg/api"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var client api.EventServiceClient

func TestMain(m *testing.M) {
	host := os.Getenv("CALENDAR_GRPC_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("CALENDAR_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	addr := fmt.Sprintf("%s:%s", host, port)

	var conn *grpc.ClientConn
	var err error

	// Retry connection because calendar service might still be starting
	for i := range 15 {
		conn, err = grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			// Verify connection is healthy
			client = api.NewEventServiceClient(conn)
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			_, err = client.Hello(ctx, &api.HelloRequest{})
			cancel()
			if err == nil {
				break
			}
		}
		fmt.Printf("Attempt %d: failed to connect to %s: %v\n", i+1, addr, err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		fmt.Printf("Could not connect to calendar service at %s: %v\n", addr, err)
		os.Exit(1)
	}

	code := m.Run()
	_ = conn.Close()
	os.Exit(code)
}

func TestCalendarCRUD(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	var eventID string
	start := time.Now().Add(time.Hour).Truncate(time.Second)
	end := start.Add(time.Hour)
	userID := "test_user"

	t.Run("create event", func(t *testing.T) {
		createReq := &api.CreateEventRequest{
			Event: &api.Event{
				Title:       "Integration Test Event",
				Description: "Test Description",
				StartAt:     timestamppb.New(start),
				EndAt:       timestamppb.New(end),
				UserId:      userID,
			},
		}

		createResp, err := client.CreateEvent(ctx, createReq)
		require.NoError(t, err)
		require.NotNil(t, createResp.Event)
		require.NotEmpty(t, createResp.Event.Id)

		eventID = createResp.Event.Id
	})

	t.Run("get event", func(t *testing.T) {
		require.NotEmpty(t, eventID)
		getResp, err := client.GetEvent(ctx, &api.GetEventRequest{Id: eventID})
		require.NoError(t, err)
		require.Equal(t, "Integration Test Event", getResp.Event.Title)
		require.Equal(t, "Test Description", getResp.Event.Description)
	})

	t.Run("update event", func(t *testing.T) {
		require.NotEmpty(t, eventID)
		updateReq := &api.UpdateEventRequest{
			Event: &api.Event{
				Id:          eventID,
				Title:       "Updated Integration Test Event",
				Description: "Updated Test Description",
				StartAt:     timestamppb.New(start),
				EndAt:       timestamppb.New(end),
				UserId:      userID,
			},
		}
		updateResp, err := client.UpdateEvent(ctx, updateReq)
		require.NoError(t, err)
		require.Equal(t, updateReq.Event.Title, updateResp.Event.Title)
	})

	t.Run("list events", func(t *testing.T) {
		require.NotEmpty(t, eventID)
		listResp, err := client.ListEvents(ctx, &api.ListEventsRequest{})
		require.NoError(t, err)
		found := false
		for _, e := range listResp.Events {
			if e.Id == eventID {
				found = true
				break
			}
		}
		require.True(t, found)
	})
}
