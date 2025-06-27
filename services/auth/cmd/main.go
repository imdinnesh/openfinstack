package main

import (
	"github.com/gin-gonic/gin"
	Logger "github.com/imdinnesh/openfinstack/packages/logger"
	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/internal/handler"
)

func main() {
	Logger.Log.Info().Msg("Starting Auth Service")
	cfg := config.Load()
	router := gin.Default()
	router.GET("/auth-test", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Auth Service is running",
		})
	})
	public := router.Group("/api/auth")
	{
		public.POST("/register", handler.Register)
		public.POST("/login", handler.Login)
		public.POST("/refresh", handler.RefreshToken)
	}

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		Logger.Log.Fatal().Err(err).Msg("Failed to start Auth Service")
	}
	Logger.Log.Info().Msg("Auth Service started successfully")
}
