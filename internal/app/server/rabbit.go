package server

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type rabbitMessage struct {
	CmdCode string `json:"cmdCode"`
	Exist   bool   `json:"exist"`
}

func formRabbitMsg(cmdCode string, successful int) (*amqp.Publishing, error) {
	msg := rabbitMessage{
		CmdCode: cmdCode,
		Exist:   successful > 0,
	}
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("error marshalling message")
	}
	return &amqp.Publishing{Body: data}, nil
}

func (s *server) rabbitPublish(ctx context.Context, publishing *amqp.Publishing) error {
	s.logger.Info("start rabbit publish")
	defer s.logger.Info("stop rabbit publish")

	ctx, cancel := context.WithTimeout(ctx, time.Duration(s.cfg.RabbitPublishTimeoutSec)*time.Second)
	defer cancel()
	err := s.rabbitChannel.PublishWithContext(
		ctx,
		s.cfg.RabbitExchangeName,
		s.cfg.RabbitKey,
		false,
		false,
		*publishing)
	if err != nil {
		return fmt.Errorf("rabbit publish error: %w", err)
	}
	return nil
}
