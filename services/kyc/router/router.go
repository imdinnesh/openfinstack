package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/auth/redis"
	"github.com/imdinnesh/openfinstack/services/kyc/config"
	"github.com/imdinnesh/openfinstack/services/kyc/internal/routes"
	"gorm.io/gorm"
)

func New(cfg *config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	// test route
	router.GET("/kyc-test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "KYC Service is running",
		})
	})

	public:=router.Group("/api/v1")
	routes.RegisterKYCRoutes(public, db, cfg, redisClient)
	return router
}
