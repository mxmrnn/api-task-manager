package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RPCClient struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRPCClient(amqpURL string) (*RPCClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq dial failed: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RPCClient{
		conn:    conn,
		channel: ch,
	}, nil
}

func (c *RPCClient) Call(ctx context.Context, queue string, payload interface{}) ([]byte, error) {
	correlationID := uuid.NewString()

	replyQueue, err := c.channel.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare reply queue: %w", err)
	}

	// Подписываемся на ответ
	msgs, err := c.channel.Consume(
		replyQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to consume reply queue: %w", err)
	}

	// Подготавливаем тело запроса
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	err = c.channel.PublishWithContext(ctx,
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: correlationID,
			ReplyTo:       replyQueue.Name,
			Body:          body,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to publish rpc request: %w", err)
	}

	select {
	case msg := <-msgs:
		if msg.CorrelationId == correlationID {
			return msg.Body, nil
		}
		return nil, fmt.Errorf("correlation id mismatch")

	case <-ctx.Done():
		return nil, fmt.Errorf("rpc call cancelled: %w", ctx.Err())

	case <-time.After(15 * time.Second):
		return nil, fmt.Errorf("rpc timeout after 15s calling %s", queue)
	}
}

func (c *RPCClient) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
