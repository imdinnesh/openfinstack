package kafka

import (
    "github.com/segmentio/kafka-go"
    "time"
)

type KafkaConfig struct {
    Brokers []string
    Topic   string
}

func NewWriter(brokers []string, topic string) *kafka.Writer {
    return kafka.NewWriter(kafka.WriterConfig{
        Brokers:      brokers,
        Topic:        topic,
        Balancer:     &kafka.LeastBytes{},
        WriteTimeout: 10 * time.Second,
    })
}

func NewReader(brokers []string, topic, groupID string) *kafka.Reader {
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers: brokers,
        Topic:   topic,
        GroupID: groupID,
    })
}
