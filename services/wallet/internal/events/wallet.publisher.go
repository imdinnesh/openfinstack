package events

import (
	"context"
	"encoding/json"

	"github.com/imdinnesh/openfinstack/packages/kafka"
)

type FundsEvent struct {
	Type   string `json:"type"`
	UserId uint   `json:"user_id"`
	Amount int64  `json:"amount"`
}

type TransactionEvent struct {
	Type   string `json:"type"`
	From         uint   `json:"from_user_id"`
	To           uint   `json:"to_user_id"`
	Amount int64  `json:"amount"`
}

type WalletEventPublisher struct {
	publisher kafka.Publisher
}

func NewWalletEventPublisher() *WalletEventPublisher {
	return &WalletEventPublisher{
		publisher: kafka.NewEventPublisher("wallet-events"),
	}
}

func (p *WalletEventPublisher) PublishAddFunds(ctx context.Context, userId uint, amount int64) error {
	event := FundsEvent{
		Type:   "WalletCredited",
		UserId: userId,
		Amount: amount,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, "wallet-events", data)
}

func (p *WalletEventPublisher) PublishDebitFunds(ctx context.Context, userId uint, amount int64) error {
	event := FundsEvent{
		Type:   "WalletDebited",
		UserId: userId,
		Amount: amount,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, "wallet-events", data)
}

func (p *WalletEventPublisher) PublishTransaction(ctx context.Context, fromUserId, toUserId uint, amount int64) error {
	event := TransactionEvent{
		Type:   "WalletTransfer",
		From:   fromUserId,
		To:     toUserId,
		Amount: amount,
	}

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, "wallet-events", data)
}
