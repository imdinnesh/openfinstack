package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/ledger/config"
	"github.com/imdinnesh/openfinstack/services/ledger/internal/routes"
	"github.com/imdinnesh/openfinstack/services/ledger/middleware"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.RequestID())

	// test route
	router.GET("/ledger-test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Ledger Service is running",
		})
	})

	public := router.Group("/api/v1")
	routes.RegisterLedgerRoutes(public, db, cfg)
	return router
}
