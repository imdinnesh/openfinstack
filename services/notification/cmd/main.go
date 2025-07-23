package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/imdinnesh/openfinstack/packages/kafka"
	consumer "github.com/imdinnesh/openfinstack/services/notifications/events"
	Logger "github.com/imdinnesh/openfinstack/packages/logger"
)

func main() {
	Logger.Log.Info().Msg("Starting Notification Service")
	dispatch := kafka.NewDispatcher()
	dispatch.RegisterHandler("user.created", consumer.HandleUserCreated)

    consumer := kafka.NewConsumer("localhost:9092", "notification-group", []string{"user.created"}, dispatch)

    ctx, cancel := context.WithCancel(context.Background())
    go func() {
        if err := consumer.Start(ctx); err != nil {
            log.Fatalf("[Kafka] Failed to start consumer: %v", err)
		} else {
			log.Println("[Kafka] Notification consumer started successfully")
        }
    }()

    // Graceful shutdown
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
    <-sigCh
    cancel()
    Logger.Log.Info().Msg("Cleaning up resources...")
}