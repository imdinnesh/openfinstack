package app

import (
	"context"
	"github.com/imdinnesh/openfinstack/packages/kafka"
	"github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/notifications/config"
	"github.com/imdinnesh/openfinstack/services/notifications/email"
	"github.com/imdinnesh/openfinstack/services/notifications/events"
)

func Run(ctx context.Context, cfg *config.Config) {
	dispatch := kafka.NewDispatcher()

	smtpSender := email.NewSMTPSender(
		cfg.SMTPUser,
		cfg.SMTPPassword,
		cfg.SMTPHost,
		cfg.SMTPPort)

	emailService := email.NewService(smtpSender)
	userHandler := consumer.NewUserCreatedHandler(emailService)
	kycHandler := consumer.NewKycStatusHandler(emailService)
	dispatch.RegisterHandler("user.created", userHandler.Handle)
	dispatch.RegisterHandler("kyc.status", kycHandler.Handle)

	consumer := kafka.NewConsumer("localhost:9092", "notification-group", []string{"user.created", "kyc.status"}, dispatch)

	ctx, cancel := context.WithCancel(ctx)

	go func() {
		if err := consumer.Start(ctx); err != nil {
			logger.Log.Fatal().Err(err).Msg("[Kafka] Consumer error")
		}
	}()

	waitForShutdown(cancel)
}
