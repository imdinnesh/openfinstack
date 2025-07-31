package main

import (
	"context"

	"github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/verifications/app"
	"github.com/imdinnesh/openfinstack/services/verifications/config"
)

func main() {
	logger.Log.Info().Msg("Starting Notification Service")
	cfg := config.Load()
	ctx := context.Background()
	app.Run(ctx,cfg)
}
