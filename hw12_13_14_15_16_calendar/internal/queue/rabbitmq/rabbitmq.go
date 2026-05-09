package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/logger"
	"github.com/Romasmi/golang-pro-course/hw12_13_14_15_calendar/internal/queue"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	logger  *logger.Logger
}

func New(url, queueName string, l *logger.Logger) (*RabbitMQ, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: ch,
		queue:   queueName,
		logger:  l,
	}, nil
}

func (r *RabbitMQ) Publish(ctx context.Context, n queue.Notification) error {
	body, err := json.Marshal(n)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	err = r.channel.PublishWithContext(ctx,
		"",      // exchange
		r.queue, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish a message: %w", err)
	}

	return nil
}

func (r *RabbitMQ) Receive(ctx context.Context) (<-chan queue.Notification, error) {
	msgs, err := r.channel.Consume(
		r.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to register a consumer: %w", err)
	}

	notifications := make(chan queue.Notification)

	go func() {
		defer close(notifications)
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				var n queue.Notification
				if err := json.Unmarshal(d.Body, &n); err != nil {
					r.logger.Error(fmt.Sprintf("failed to unmarshal notification: %v, body: %s", err, string(d.Body)))
					continue
				}
				notifications <- n
			}
		}
	}()

	return notifications, nil
}

func (r *RabbitMQ) Close() error {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		return r.conn.Close()
	}
	return nil
}
