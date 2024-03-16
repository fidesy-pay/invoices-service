package consumers

import (
	"context"
	"fmt"
	"github.com/fidesy-pay/invoices-service/internal/config"
	"github.com/fidesy/sdk/common/kafka"
	"github.com/fidesy/sdk/common/logger"
)

type MessageHandler interface {
	Consume(ctx context.Context, msg []byte) error
}

func RegisterConsumer(ctx context.Context, handler MessageHandler, topic string) error {
	kafkaConsumer, err := kafka.NewConsumer(ctx, config.Get(config.KafkaBrokers).([]string), topic)
	if err != nil {
		return fmt.Errorf("kafka.NewConsumer: %w", err)
	}

	go consume(ctx, handler, kafkaConsumer)

	return nil
}

func consume(ctx context.Context, handler MessageHandler, consumer *kafka.Consumer) {
	defer func() {
		err := consumer.Close()
		if err != nil {
			logger.Fatalf("kafkaConsumer.Close: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-consumer.Consume():
			err := handler.Consume(ctx, msg.Value)
			if err != nil {
				logger.Errorf("handler.Consume: %v", err)
			}
		}
	}
}
