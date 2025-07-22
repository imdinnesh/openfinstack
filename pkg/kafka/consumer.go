package kafka

import (
    "context"
    "log"

    "github.com/segmentio/kafka-go"
)

type HandlerFunc func(key, value []byte)

func Consume(ctx context.Context, reader *kafka.Reader, handler HandlerFunc) {
    for {
        m, err := reader.ReadMessage(ctx)
        if err != nil {
            log.Printf("Error reading message: %v", err)
            continue
        }
        handler(m.Key, m.Value)
    }
}
