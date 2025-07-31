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

type KYCEventPublisher struct {
    publisher kafka.Publisher
}

func NewKYCEventPublisher() *KYCEventPublisher {
    return &KYCEventPublisher{
        publisher: kafka.NewEventPublisher("kyc.submitted"),
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
