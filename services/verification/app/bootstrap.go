package app

import (
	"context"

	"github.com/imdinnesh/openfinstack/packages/kafka"
	"github.com/imdinnesh/openfinstack/packages/logger"
	clients "github.com/imdinnesh/openfinstack/services/verifications/client"
	"github.com/imdinnesh/openfinstack/services/verifications/config"
	consumer "github.com/imdinnesh/openfinstack/services/verifications/events"
	repository "github.com/imdinnesh/openfinstack/services/verifications/repo"
	"github.com/imdinnesh/openfinstack/services/verifications/service"
	"github.com/imdinnesh/openfinstack/services/verifications/verifier/provider"
	"gorm.io/gorm"
)

func Run(ctx context.Context, cfg *config.Config,db *gorm.DB) {
	dispatch := kafka.NewDispatcher()
	verifier := provider.NewVerifier(cfg)
	kycRepo := repository.NewKYCRepository(db)
	kycClient := clients.NewClient(cfg.KYCBaseURL)
	verifierService := service.NewService(verifier, kycRepo, kycClient)
	kycHandler := consumer.NewKYCHandler(verifierService)

	dispatch.RegisterHandler("kyc.submitted", kycHandler.Handle)

	consumer := kafka.NewConsumer("localhost:9092", "verification-group", []string{"kyc.submitted"}, dispatch)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.Log.Fatal().Err(err).Msg("[Kafka] Consumer error")
		}
	}()

	waitForShutdown(cancel)
}
