package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/services/auth/config"
	"github.com/imdinnesh/openfinstack/services/auth/internal/handler"
)

func New(cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// test route
	router.GET("/auth-test", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Auth Service is running",
		})
	})

	public := router.Group("/api/auth")
	public.POST("/register", handler.Register)
	public.POST("/login", handler.Login)
	public.POST("/refresh", handler.RefreshToken)
	return router
}
