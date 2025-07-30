package oauth
import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterOAuthRoutes(r *gin.RouterGroup, db *gorm.DB) {
	service := NewService(db)
	handler := NewHandler(service)
	oauth := r.Group("/oauth")
	oauth.GET("/google/login", handler.GoogleLogin)
	oauth.GET("/google/callback", handler.GoogleCallback)
}