package rabbitmq

import (
	"context"
	"testing"
	"time"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/queue"
	"github.com/stretchr/testify/require"
)

// This test requires a running RabbitMQ instance.
// If it's not available, the test will be skipped.
func TestRabbitMQ(t *testing.T) {
	url := "amqp://rabbit:password@localhost:5672/"
	queueName := "test_queue"

	rmq, err := New(url, queueName)
	if err != nil {
		t.Skip("RabbitMQ is not available")
		return
	}
	defer rmq.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	n := queue.Notification{
		EventID: "123",
		Title:   "Test",
		UserID:  "user1",
		StartAt: time.Now().Round(time.Second),
	}

	err = rmq.Publish(ctx, n)
	require.NoError(t, err)

	msgs, err := rmq.Receive(ctx)
	require.NoError(t, err)

	select {
	case got := <-msgs:
		require.Equal(t, n.EventID, got.EventID)
		require.Equal(t, n.Title, got.Title)
		require.Equal(t, n.UserID, got.UserID)
		require.True(t, n.StartAt.Equal(got.StartAt))
	case <-ctx.Done():
		t.Fatal("timeout waiting for message")
	}
}
