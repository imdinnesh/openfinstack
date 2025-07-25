package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/packages/middleware"
	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/internal/events"
	"github.com/imdinnesh/openfinstack/services/auth/internal/handler"
	"github.com/imdinnesh/openfinstack/services/auth/internal/repository"
	"github.com/imdinnesh/openfinstack/services/auth/internal/service"
	"github.com/imdinnesh/openfinstack/services/auth/redis"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.RouterGroup, db *gorm.DB, cfg *config.Config, redisClient *redis.Client) {
	userRepo := repository.NewUserRepository(db)
	publisher := events.NewUserEventPublisher()
	authSvc := service.NewAuthService(userRepo, cfg, redisClient, publisher)
	authHandler := handler.NewAuthHandler(authSvc)

	middleware.New(cfg.JWTSecret, redisClient)

	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)
	auth.POST("/logout",authHandler.Logout)
	auth.GET("/profile",authHandler.Profile)
}
