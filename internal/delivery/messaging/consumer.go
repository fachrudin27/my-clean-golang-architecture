package messaging

import (
	"fmt"
	"my-clean-architecture-template/pkg/rabbitmq"
)

func InitConsumer(conn *rabbitmq.Connection) {
	chanName := []string{
		"detection_ai_response",
		"detection_vital_chubb_response",
	}

	messageChannels := rabbitmq.InitializeMessageChannels(chanName)

	for queueName, messageChan := range messageChannels {
		go func(queueName string, messageChan chan []byte) {
			conn.ConsumeQueue(queueName, messageChan, conn.Stop)
		}(queueName, messageChan)
	}

	go processMessages(messageChannels, conn.Stop)
}

func processMessages(messageChannels map[string]chan []byte, shutdown chan struct{}) {
	for {
		select {
		case msg := <-messageChannels["detection_ai_response"]:
			LoginConsume(msg)
			// fmt.Println(string(msg))
		case msg := <-messageChannels["detection_vital_chubb_response"]:
			fmt.Println(string(msg))
		case <-shutdown:
			fmt.Println("Stopping message processor...")
			return
		}
	}
}
