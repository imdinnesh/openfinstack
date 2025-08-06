package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/wallet/config"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/events"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/handler"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/repository"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/service"
	"gorm.io/gorm"
)


func RegisterWalletRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config) {
	publisher := events.NewWalletEventPublisher()
	walletRepo := repository.New(db, publisher)
	walletSvc := service.New(walletRepo)
	walletHandler := handler.New(walletSvc)

	wallet := r.Group("/wallet")
	wallet.POST("/wallet", walletHandler.CreateWallet)
	wallet.GET("/:userID", walletHandler.GetWallet)
	wallet.POST("/:userID/credit", walletHandler.AddFunds)
	wallet.POST("/:userID/debit", walletHandler.WithdrawFunds)
	wallet.POST("/transfer", walletHandler.Transfer)
	wallet.GET("/:userID/transactions", walletHandler.GetTransactions)

}