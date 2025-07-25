package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	constants "github.com/imdinnesh/openfinstack/packages/config"
	Logger "github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/wallet/config"
	"github.com/imdinnesh/openfinstack/services/wallet/db"
	"github.com/imdinnesh/openfinstack/services/wallet/router"
)

func main() {
	Logger.Log.Info().Msg("Starting Wallet Service")
	cfg := config.Load()
	db := db.InitDB(cfg)
	
	// Set up router (Gin)
	Router := router.New(cfg, db)

	// Create custom HTTP server
	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: Router,
	}

	// Signal handling for graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		Logger.Log.Info().Msgf("Wallet Service listening on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			Logger.Log.Error().Err(err).Msg("Server failed to start")
		} else {
			Logger.Log.Info().Msg("Server started successfully")
		}
	}()

	<-done // Wait for signal

	Logger.Log.Info().Msg("Received shutdown signal, shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), constants.ShutdownTimeout)
	defer cancel()

	// Graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		Logger.Log.Error().Err(err).Msg("Server forced to shutdown")
	} else {
		Logger.Log.Info().Msg("Server gracefully stopped")
	}

	Logger.Log.Info().Msg("Cleaning up resources...")
}
