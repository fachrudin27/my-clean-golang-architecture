package rabbitmq

import (
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (conn *Connection) QueueDeclare(queueName string) error {
	_, err := conn.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue %v", err)
	}

	return nil
}

func (conn *Connection) QueuePublish(queueName string, body []byte) error {
	err := conn.Channel.Publish(
		"",
		queueName,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to queue %s: %v", queueName, err)
	}
	return nil
}

// Notify -.
func (s *Connection) Notify() <-chan error {
	return s.error
}

// Shutdown -.
func (conn *Connection) Shutdown() error {
	close(conn.Stop)
	time.Sleep(5 * time.Second)

	if conn.Connection != nil {
		if err := conn.Connection.Close(); err != nil {
			return fmt.Errorf("failed to close RabbitMQ connection: %w", err)
		}
	}

	conn.Logger.Info("RabbitMQ connection shutdown successfully")
	return nil
}

func (conn *Connection) Consumer(queueName string) {

	q, err := conn.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Logger.Error("Failed to declare queue name : %s, err : %v", queueName, err)
	}

	msgs, err := conn.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Logger.Error("Failed to consume messages: %v", err)
	}

	conn.Logger.Info("Listening for messages asynchronously...")

	go func() {
		for msg := range msgs {
			conn.Logger.Info("Received async message: %s", string(msg.Body))
		}
	}()
}

func InitializeMessageChannels(keys []string) map[string]chan []byte {
	messageChannels := make(map[string]chan []byte)
	for _, key := range keys {
		messageChannels[key] = make(chan []byte)
	}
	return messageChannels
}

func (conn *Connection) ConsumeQueue(queueName string, messages chan []byte, shutdown chan struct{}) {
	q, err := conn.Channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		conn.Logger.Error("Failed to declare queue %s: %v\n", queueName, err)
	}

	msgs, err := conn.Channel.Consume(
		q.Name,
		"",
		false, // auto-ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		conn.Logger.Error("Failed to start consuming from queue %s: %v\n", queueName, err)
	}

	for {
		select {
		case msg := <-msgs:
			messages <- msg.Body
			if err := msg.Ack(false); err != nil {
				conn.Logger.Info("Failed to ack message from queue %s: %v\n", queueName, err)
			}
		case <-shutdown:
			conn.Logger.Info("Stopping consumer for queue %s...\n", queueName)
			return
		}
	}
}
