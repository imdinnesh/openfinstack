package events

import (
	"encoding/json"
	"github.com/imdinnesh/openfinstack/services/notifications/email"
)

type UserCreatedEvent struct {
    ID    uint   `json:"id"`
    Email string `json:"email"`
}

type UserCreatedHandler struct {
    EmailService *email.Service
}

func NewUserCreatedHandler(es *email.Service) *UserCreatedHandler {
    return &UserCreatedHandler{EmailService: es}
}

func (h *UserCreatedHandler) Handle(key, value []byte) error {
    var event UserCreatedEvent
    if err := json.Unmarshal(value, &event); err != nil {
        return err
    }

    template := email.OnboardingEmail{
        UserID: event.ID,
        Email:  event.Email,
    }

    return h.EmailService.Send(event.Email, template)
}
