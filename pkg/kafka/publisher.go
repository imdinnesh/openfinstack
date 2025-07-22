package kafka

import (
    "context"
    "github.com/segmentio/kafka-go"
)

type EventPublisher struct {
    writer *kafka.Writer
}

func NewEventPublisher(topic string) Publisher {
    return &EventPublisher{
        writer: getOrCreateWriter(topic),
    }
}

func (p *EventPublisher) Publish(ctx context.Context, key string, value []byte) error {
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(key),
        Value: value,
    })
}
