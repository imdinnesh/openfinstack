package events

import (
	"encoding/json"
	"fmt"
	
	"github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/service"
	"github.com/imdinnesh/openfinstack/services/ledger/models"
)

// WalletEvent represents the generic structure of an event coming from the Wallet service.
// The actual payload is determined by the `Type` field.
type WalletEvent struct {
	Type     string `json:"type"`
	UserID   uint   `json:"user_id"`
	Amount   int64  `json:"amount"`
	From     uint   `json:"from_user_id"`
	To       uint   `json:"to_user_id"`
}

// WalletEventHandler holds the dependencies for handling wallet events.
type WalletEventHandler struct {
	ledgerService service.LedgerService
}

// NewWalletEventHandler creates a new handler instance.
func NewWalletEventHandler(ledgerService service.LedgerService) *WalletEventHandler {
	return &WalletEventHandler{
		ledgerService: ledgerService,
	}
}

// Handle processes the incoming Kafka message from the wallet events topic.
// It parses the event and creates a new ledger transaction.
func (h *WalletEventHandler) Handle(key, value []byte) error {
	var event WalletEvent
	if err := json.Unmarshal(value, &event); err != nil {
		logger.Log.Error().Err(err).Msg("Failed to unmarshal wallet event JSON")
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	logger.Log.Info().Msgf("Received wallet event of type: %s", event.Type)

	var req service.CreateTxnRequest
	
	// Create a new Ledger transaction based on the event type.
	switch event.Type {
	case "WalletCredited":
		req = service.CreateTxnRequest{
			ExternalRef: fmt.Sprintf("wallet_credit_%d_%d", event.UserID, event.Amount), // Simple idempotency key
			Type:        models.TxnWalletTopup,
			Description: fmt.Sprintf("Credit to user wallet %d", event.UserID),
			Currency:    "USD", // Hardcoded currency for now
			Entries: []service.CreateEntry{
				{
					AccountID: fmt.Sprintf("user_wallet_%d", event.UserID),
					EntryType: models.DEBIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
				{
					AccountID: "internal_funds",
					EntryType: models.CREDIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
			},
		}
	case "WalletDebited":
		req = service.CreateTxnRequest{
			ExternalRef: fmt.Sprintf("wallet_debit_%d_%d", event.UserID, event.Amount),
			Type:        models.TxnWithdrawal,
			Description: fmt.Sprintf("Debit from user wallet %d", event.UserID),
			Currency:    "USD",
			Entries: []service.CreateEntry{
				{
					AccountID: fmt.Sprintf("user_wallet_%d", event.UserID),
					EntryType: models.CREDIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
				{
					AccountID: "internal_funds",
					EntryType: models.DEBIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
			},
		}
	case "WalletTransfer":
		req = service.CreateTxnRequest{
			ExternalRef: fmt.Sprintf("wallet_transfer_%d_%d_%d", event.From, event.To, event.Amount),
			Type:        models.TxnWalletTransfer,
			Description: fmt.Sprintf("Transfer from user %d to user %d", event.From, event.To),
			Currency:    "USD",
			Entries: []service.CreateEntry{
				{
					AccountID: fmt.Sprintf("user_wallet_%d", event.From),
					EntryType: models.DEBIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
				{
					AccountID: fmt.Sprintf("user_wallet_%d", event.To),
					EntryType: models.CREDIT,
					Amount:    event.Amount,
					Meta:      nil,
				},
			},
		}
	default:
		logger.Log.Warn().Msgf("Unknown wallet event type: %s", event.Type)
		return nil
	}

	// Submit the new transaction to the Ledger service.
	_, err := h.ledgerService.CreateAndPost(req)
	if err != nil {
		logger.Log.Error().Err(err).Msgf("Failed to create ledger transaction from wallet event for user %d", event.UserID)
		return fmt.Errorf("failed to create ledger transaction: %w", err)
	}

	logger.Log.Info().Msgf("Successfully created ledger transaction for wallet event: %s", event.Type)
	return nil
}