package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/imdinnesh/openfinstack/packages/logger"
)

func waitForShutdown(cancelFunc func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigCh

	logger.Log.Info().Msgf("Shutting down due to signal: %s", sig.String())
	cancelFunc()
}
