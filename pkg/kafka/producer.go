package kafka

import (
    "context"
    "github.com/segmentio/kafka-go"
)

type Producer struct {
    Writer *kafka.Writer
}

func NewProducer(writer *kafka.Writer) *Producer {
    return &Producer{Writer: writer}
}

func (p *Producer) Publish(ctx context.Context, key string, value []byte) error {
    return p.Writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(key),
        Value: value,
    })
}
