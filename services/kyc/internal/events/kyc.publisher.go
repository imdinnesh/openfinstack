package events

import (
	"context"
	"encoding/json"

	"github.com/imdinnesh/openfinstack/packages/kafka"
)

type KYCDocumentSubmittedEvent struct {
    KYCID       uint   `json:"kyc_id"`
    DocumentType string `json:"document_type"`
	DocumentURL  string `json:"document_url"`
}

type KycStatusEvent struct {
    UserID uint
    Email  string
    Status string // e.g. "Approved", "Rejected", "Pending"
    Reason string // Optional: reason for rejection or additional info
}

type KYCEventPublisher struct {
    publisher kafka.Publisher
    publisher2 kafka.Publisher
}

func NewKYCEventPublisher() *KYCEventPublisher {
    return &KYCEventPublisher{
        publisher: kafka.NewEventPublisher("kyc.submitted"),
        publisher2: kafka.NewEventPublisher("kyc.status"),
    }
}

func (p *KYCEventPublisher) PublishKYCDocumentSubmitted(ctx context.Context, kycID uint, documentType, documentURL string) error {
    event := KYCDocumentSubmittedEvent{
        KYCID:       kycID,
        DocumentType: documentType,
        DocumentURL:  documentURL,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.publisher.Publish(ctx, "kyc.submitted", data)
}

func (p *KYCEventPublisher) PublishKYCStatus(ctx context.Context, userID uint, email, status, reason string) error {
    event := KycStatusEvent{
        UserID: userID,
        Email:  email,
        Status: status,
        Reason: reason,
    }

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.publisher2.Publish(ctx, "kyc.status", data)
}
