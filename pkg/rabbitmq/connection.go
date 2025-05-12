package rabbitmq

import (
	"fmt"
	"log"
	"my-clean-architecture-template/pkg/logger"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	URL      string
	WaitTime time.Duration
	Attempts int
	Stop     chan struct{}

	Logger    logger.Interface
	QueueName string
}

type Connection struct {
	Config
	Connection *amqp.Connection
	Channel    *amqp.Channel
	error      chan error
}

func New(cfg Config) (*Connection, error) {
	conn := &Connection{
		Config: cfg,
	}

	err := conn.AttemptConnect()
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (conn *Connection) AttemptConnect() error {
	var err error
	for i := conn.Attempts; i > 0; i-- {
		if err = conn.connect(); err == nil {
			break
		}

		log.Printf("RabbitMQ is trying to connect, attempts left: %d", i)
		time.Sleep(conn.WaitTime)
	}

	if err != nil {
		return fmt.Errorf("failed to attempt connect rabbitmq, AttemptConnect: %w", err)
	}

	return nil
}

func (conn *Connection) connect() error {
	var err error

	conn.Connection, err = amqp.Dial(conn.Config.URL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	conn.Channel, err = conn.Connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	return nil
}
