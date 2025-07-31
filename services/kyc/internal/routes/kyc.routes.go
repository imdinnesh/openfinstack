package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/handler"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/repository"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/service"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/verifier/provider"
	"gorm.io/gorm"
)

func RegisterKYCRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	kycRepo := repository.NewKYCRepository(db)
	var verifier verifier.Verifier

	switch cfg.KYCVerifier {
	case "idfy":
		verifier = provider.NewIDfyVerifier(cfg.IDFYApiKey, cfg.IDFYBaseURL)
	default:
		verifier = provider.NewMockVerifier()
	}
	kycSvc := service.NewKYCService(kycRepo, verifier)
	kycHandler := handler.NewKYCHandler(kycSvc)

	kyc := r.Group("/kyc")
	kyc.POST("/submit", kycHandler.SubmitKYC)
	kyc.GET("/user", kycHandler.GetUserKYC)
	kyc.GET("/status", kycHandler.GetKYCStatusByUserID)
	kycAdmin := r.Group("kyc-admin")
	kycAdmin.GET("/pending", kycHandler.ListPending)
	kycAdmin.POST("/verify/:id", kycHandler.VerifyKYC)
}
