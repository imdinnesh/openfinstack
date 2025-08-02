package routes

import (
	"github.com/gin-gonic/gin"
	clients "github.com/imdinnesh/openfinstack/services/kyc/client"
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/events"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/handler"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/repository"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/service"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier/provider"
	"gorm.io/gorm"
)

func RegisterKYCRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	kycRepo := repository.NewKYCRepository(db)
	verifier := provider.NewVerifier(cfg)
	publisher := events.NewKYCEventPublisher()
	userClient := clients.NewClient(cfg.UserServiceURL)
	kycSvc := service.NewKYCService(kycRepo, verifier, publisher, userClient)
	kycHandler := handler.NewKYCHandler(kycSvc)

	kyc := r.Group("/kyc")
	kyc.POST("/submit", kycHandler.SubmitKYC)
	kyc.GET("/user", kycHandler.GetUserKYC)
	kyc.GET("/status", kycHandler.GetKYCStatusByUserID)
	kyc.POST("/update/:id", kycHandler.UpdateKYCStatus)
	kycAdmin := r.Group("kyc-admin")
	kycAdmin.GET("/pending", kycHandler.ListPending)
	kycAdmin.POST("/verify/:id", kycHandler.VerifyKYC)
}
