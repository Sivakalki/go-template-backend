package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers[] string, topic string) *Producer{
	return &Producer{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic: topic,
			Balancer: &kafka.LeastBytes{},
		}),
	}
}


func (p* Producer) Publish(ctx context.Context, key , value string) error {
	msg:= kafka.Message{
		Key: []byte(key),
		Value: []byte(value),
	}

	return p.writer.WriteMessages(ctx, msg)
}

func (p* Producer) Close() error{
	return p.writer.Close()
}