package discovery

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// ProxyHandler forwards the incoming request to the target service.
func ProxyHandler(baseURL string, servicePath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		finalPath := servicePath
		for _, param := range c.Params {
			finalPath = strings.ReplaceAll(finalPath, ":"+param.Key, param.Value)
		}

		targetURL := baseURL + finalPath

		var body io.Reader = nil
		if c.Request.Body != nil {
			buf, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
				return
			}
			body = bytes.NewReader(buf)
			c.Request.Body = io.NopCloser(bytes.NewReader(buf))
		}

		req, err := http.NewRequest(c.Request.Method, targetURL, body)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
			return
		}

		for key, vals := range c.Request.Header {
			for _, val := range vals {
				req.Header.Add(key, val)
			}
		}

		client := http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Failed to reach target service"})
			return
		}
		defer resp.Body.Close()

		c.Status(resp.StatusCode)
		for key, vals := range resp.Header {
			for _, val := range vals {
				c.Header(key, val)
			}
		}
		io.Copy(c.Writer, resp.Body)
	}
}
