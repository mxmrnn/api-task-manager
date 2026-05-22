package rpc

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RPCServer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewServer(amqpURL string) (*RPCServer, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &RPCServer{conn: conn, channel: ch}, nil
}

func (s *RPCServer) Start(queueName string, handler func([]byte) ([]byte, error)) error {
	// Объявляем очередь
	_, err := s.channel.QueueDeclare(
		queueName, // task-service.rpc
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	msgs, err := s.channel.Consume(
		queueName,
		"",    // consumer tag
		false, // autoAck = false (вручную)
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	go func() {
		for msg := range msgs {
			response, err := handler(msg.Body)
			if err != nil {
				log.Printf("[RPC] Handler erroror correlation %s: %v", msg.CorrelationId, err)
				errorResp := map[string]string{
					"error":   "internal_error",
					"message": err.Error(),
				}
				response, _ = json.Marshal(errorResp)
			}

			// Отправляем ответ
			s.channel.Publish(
				"",          // exchange
				msg.ReplyTo, // куда отвечать
				false,
				false,
				amqp.Publishing{
					ContentType:   "application/json",
					CorrelationId: msg.CorrelationId,
					Body:          response,
				},
			)
			if err != nil {
				log.Printf("[RPC] Failed to publish response: %v", err)
			}

			err = msg.Ack(false)
			if err != nil {
				log.Printf("[RPC] Failed to ack message: %v", err)
			}
		}
	}()
	log.Printf("RPC Server started on queue: %s", queueName)
	return nil
}

func (s *RPCServer) Close() error {
	if s.channel != nil {
		s.channel.Close()
	}
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
