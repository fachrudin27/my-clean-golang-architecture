package messaging

import (
	"encoding/json"
	"my-clean-architecture-template/internal/model"
)

var (
	LoginQueue = "detection_ai_response"
)

func LoginProducer(payload model.LoginUserRequest) error {
	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = Prod.conn.QueueDeclare(LoginQueue)
	if err != nil {
		return err
	}

	err = Prod.conn.QueuePublish(LoginQueue, value)
	if err != nil {
		return err
	}

	return nil
}
