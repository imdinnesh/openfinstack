package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	constants "github.com/imdinnesh/openfinstack/packages/config"
	"github.com/imdinnesh/openfinstack/packages/kafka"
	Logger "github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/ledger/config"
	"github.com/imdinnesh/openfinstack/services/ledger/db"
	events "github.com/imdinnesh/openfinstack/services/ledger/event"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/repository"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/service"
	"github.com/imdinnesh/openfinstack/services/ledger/router"
)

func main() {
	Logger.Log.Info().Msg("Starting Ledger Service")
	cfg := config.Load()
	database, err := db.Connect(cfg.DBUrl)
	if err != nil {
		Logger.Log.Error().Err(err).Msg("Failed to connect to database")
		return
	}
	db.Migrate(database)

	// Initialize the Ledger service dependencies
	ledgerRepo := repository.NewLedgerRepository(database)
	ledgerService := service.NewLedgerService(database, ledgerRepo)
	
	// Set up Kafka consumer
	dispatcher := kafka.NewDispatcher()
	walletEventHandler := events.NewWalletEventHandler(ledgerService)
	dispatcher.RegisterHandler("wallet-events", walletEventHandler.Handle)
	
	consumer := kafka.NewConsumer("localhost:9092", "ledger-group", []string{"wallet-events"}, dispatcher)
	
	// Start the consumer in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			Logger.Log.Fatal().Err(err).Msg("[Kafka] Consumer error")
		}
	}()
	
	// Set up router (Gin)
	Router := router.New(cfg, database)

	// Create custom HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: Router,
	}

	// Signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		Logger.Log.Info().Msgf("Ledger Service listening on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Log.Error().Err(err).Msg("Server failed to start")
		} else {
			Logger.Log.Info().Msg("Server started successfully")
		}
	}()

	<-done // Wait for signal

	Logger.Log.Info().Msg("Received shutdown signal, shutting down server...")

	ctx, srvCancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer srvCancel()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		Logger.Log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		Logger.Log.Info().Msg("Server gracefully stopped")
	}

	Logger.Log.Info().Msg("Cleaning up resources...")
}