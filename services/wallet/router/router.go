package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/wallet/config"
	"github.com/imdinnesh/openfinstack/services/wallet/internal/routes"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// test route
	router.GET("/wallet-test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Wallet Service is running",
		})
	})

	public := router.Group("/api/v1")
	routes.RegisterWalletRoutes(public, db, cfg)
	return router
}
