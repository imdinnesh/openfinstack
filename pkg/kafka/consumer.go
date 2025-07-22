package kafka

import (
    "context"
    "log"

    "github.com/segmentio/kafka-go"
)

type Consumer struct {
    broker   string
    groupID  string
    topics   []string
    dispatch *Dispatcher
}

func NewConsumer(broker, groupID string, topics []string, dispatch *Dispatcher) *Consumer {
    return &Consumer{
        broker:   broker,
        groupID:  groupID,
        topics:   topics,
        dispatch: dispatch,
    }
}

func (c *Consumer) Start(ctx context.Context) error {
    for _, topic := range c.topics {
        go func(topic string) {
            reader := kafka.NewReader(kafka.ReaderConfig{
                Brokers: c.brokerList(),
                GroupID: c.groupID,
                Topic:   topic,
            })

            log.Println("[Kafka] Consumer started for topic:", topic)

            for {
                m, err := reader.ReadMessage(ctx)
                if err != nil {
                    log.Printf("[Kafka] Error reading message: %v\n", err)
                    continue
                }

                handler, ok := c.dispatch.GetHandler(m.Topic)
                if !ok {
                    log.Printf("[Kafka] No handler registered for topic: %s\n", m.Topic)
                    continue
                }

                if err := handler(m.Key, m.Value); err != nil {
                    log.Printf("[Kafka] Handler error for topic %s: %v\n", m.Topic, err)
                }
            }
        }(topic)
    }

    // Block until context is done
    <-ctx.Done()
    return ctx.Err()
}

func (c *Consumer) brokerList() []string {
    return []string{c.broker}
}
