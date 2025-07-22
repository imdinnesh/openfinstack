package kafka

import (
	"sync"
	"time"

	constants "github.com/imdinnesh/openfinstack/packages/config"
	"github.com/segmentio/kafka-go"
)

var (
    writerCache = make(map[string]*kafka.Writer)
    writerLock  sync.Mutex
)

func getOrCreateWriter(topic string) *kafka.Writer {
    writerLock.Lock()
    defer writerLock.Unlock()

    if w, exists := writerCache[topic]; exists {
        return w
    }

    writer := &kafka.Writer{
        Addr:         kafka.TCP(constants.Brokers...),
        Topic:        topic,
        Balancer:     &kafka.LeastBytes{},
        WriteTimeout: 10 * time.Second,
    }

    writerCache[topic] = writer
    return writer
}
