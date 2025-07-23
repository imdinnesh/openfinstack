package consumer

import (
    "encoding/json"
    "log"
)

type UserCreatedEvent struct {
    ID    uint `json:"id"`
    Email string `json:"email"`
}

func HandleUserCreated(key, value []byte) error {
    var event UserCreatedEvent
    if err := json.Unmarshal(value, &event); err != nil {
        return err
    }

    log.Printf("[Notification] Send welcome email to %s (User ID: %d)\n", event.Email, event.ID)
    return nil
}
