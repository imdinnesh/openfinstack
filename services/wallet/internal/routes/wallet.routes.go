package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/wallet/config"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/handler"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/repository"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/service"
	"gorm.io/gorm"
)


func RegisterWalletRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	walletRepo := repository.New(db)
	walletSvc := service.New(walletRepo)
	walletHandler := handler.New(walletSvc)

	wallet := r.Group("/wallet")
	wallet.GET("/:userId", walletHandler.GetWallet)
	wallet.POST("/transfer", walletHandler.Transfer)
	wallet.GET("/transactions/:userId", walletHandler.Transactions)

}