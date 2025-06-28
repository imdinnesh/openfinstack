package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/internal/routes"
	"github.com/imdinnesh/openfinstack/services/auth/redis"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	// test route
	router.GET("/auth-test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Auth Service is running",
		})
	})

	public := router.Group("/api/v1")
	routes.RegisterAuthRoutes(public, db, cfg, redisClient)

	return router
}
