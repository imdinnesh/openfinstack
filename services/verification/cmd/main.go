package main

import (
	"context"

	"github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/verifications/app"
	"github.com/imdinnesh/openfinstack/services/verifications/config"
	"github.com/imdinnesh/openfinstack/services/verifications/db"
)

func main() {
	logger.Log.Info().Msg("Starting Notification Service")
	cfg := config.Load()
	db := db.InitDB(cfg)
	ctx := context.Background()
	app.Run(ctx,cfg,db)
}
