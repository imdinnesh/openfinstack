package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/handler"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/repository"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/service"
	"gorm.io/gorm"
)


func RegisterKYCRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	kycRepo := repository.NewKYCRepository(db)	
	kycSvc := service.NewKYCService(kycRepo)
	kycHandler := handler.NewKYCHandler(kycSvc)


	kyc := r.Group("/kyc")
	kyc.POST("/submit",kycHandler.SubmitKYC)
	kyc.GET("/user", kycHandler.GetUserKYC)
	kycAdmin:=r.Group("kyc-admin")
	kycAdmin.GET("/pending", kycHandler.ListPending)
	kycAdmin.POST("/verify/:id", kycHandler.VerifyKYC)
}