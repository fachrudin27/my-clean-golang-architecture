package messaging

import (
	"my-clean-architecture-template/pkg/rabbitmq"
)

var (
	Prod *Producer
)

type Producer struct {
	conn *rabbitmq.Connection
}

func InitProducer(conn *rabbitmq.Connection) *Producer {
	Prod = &Producer{
		conn: conn,
	}

	return Prod
}
