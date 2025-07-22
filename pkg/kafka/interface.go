package kafka

import "context"

// Publisher interface for Kafka message publishing
type Publisher interface {
    Publish(ctx context.Context, key string, value []byte) error
}

// MessageHandler is a function that handles Kafka messages
type MessageHandler func(key, value []byte) error
