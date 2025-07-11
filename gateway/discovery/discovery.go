package discovery

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ProxyHandler(baseURL, servicePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetURL := baseURL + servicePath

		req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header = c.Request.Header

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}
