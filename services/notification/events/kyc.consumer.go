package events

import (
	"encoding/json"
	"github.com/imdinnesh/openfinstack/services/notifications/email"
)

type KycStatusEvent struct {
    UserID uint
    Email  string
    Status string // e.g. "Approved", "Rejected", "Pending"
    Reason string // Optional: reason for rejection or additional info
}

type KycStatusHandler struct {
    EmailService *email.Service
}

func NewKycStatusHandler(es *email.Service) *KycStatusHandler {
    return &KycStatusHandler{EmailService: es}
}

func (h *KycStatusHandler) Handle(key, value []byte) error {
    var event KycStatusEvent
    if err := json.Unmarshal(value, &event); err != nil {
        return err
    }

    template := email.KYCStatusEmail{
		UserID: event.UserID,
		Email:  event.Email,
		Status: event.Status,
		Reason: event.Reason,
	}

    return h.EmailService.Send(event.Email, template)
}
