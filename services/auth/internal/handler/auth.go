package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "register"})
}

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "login"})
}

func RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "refresh token"})
}
