package consumer

import (
	"encoding/json"

	"github.com/imdinnesh/openfinstack/services/verifications/service"
)

type KYCDocumentSubmittedEvent struct {
    KYCID       uint   `json:"kyc_id"`
    DocumentType string `json:"document_type"`
	DocumentURL  string `json:"document_url"`
}

type KYCDocumentSubmittedHandler struct {
    VerifierService *service.Service
}

func NewKYCHandler(vs *service.Service) *KYCDocumentSubmittedHandler {
    return &KYCDocumentSubmittedHandler{VerifierService: vs}
}

func (h *KYCDocumentSubmittedHandler) Handle(key, value []byte) error {
    var event *KYCDocumentSubmittedEvent
    if err := json.Unmarshal(value, &event); err != nil {
        return err
    }

    kycDocument := &service.KYCDocumentSubmittedEvent{
        KYCID:       event.KYCID,
        DocumentType: event.DocumentType,
        DocumentURL:  event.DocumentURL,
    }

    return h.VerifierService.VerifyKYC(kycDocument)
}
