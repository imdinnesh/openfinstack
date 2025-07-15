package discovery

import (
	"bytes"
	"fmt"
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
      buf, _ := io.ReadAll(c.Request.Body)
      body = bytes.NewReader(buf)
      c.Request.Body = io.NopCloser(bytes.NewReader(buf))
    }

    req, err := http.NewRequest(c.Request.Method, targetURL, body)
    if err != nil {
      c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create proxy request"})
      return
    }

    // ➜ Copy original headers
    for key, vals := range c.Request.Header {
      for _, val := range vals {
        req.Header.Add(key, val)
      }
    }

    // ✅ ➜ Add trusted identity headers if present
    if userID, exists := c.Get("user_id"); exists {
      req.Header.Set("X-User-ID", fmt.Sprintf("%v", userID))
    }

    req.URL.RawQuery = c.Request.URL.RawQuery

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

