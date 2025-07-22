package events

import (
	"context"
	"encoding/json"

	"github.com/imdinnesh/openfinstack/packages/kafka"
)

type UserCreatedEvent struct {
    ID    uint   `json:"id"`
    Email string `json:"email"`
}

type UserEventPublisher struct {
    publisher kafka.Publisher
}

func NewUserEventPublisher() *UserEventPublisher {
    return &UserEventPublisher{
        publisher: kafka.NewEventPublisher("user.created"),
    }
}

func (p *UserEventPublisher) PublishUserCreated(ctx context.Context, id uint, email string) error {
    event := UserCreatedEvent{
        ID:    id,
        Email: email,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.publisher.Publish(ctx, "user.created", data)
}
