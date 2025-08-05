package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/clients"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/gateway/discovery"
	"github.com/imdinnesh/openfinstack/gateway/middleware"
	"github.com/imdinnesh/openfinstack/packages/redis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(
	cfg *config.Config,
	cfgEnvs *config.ConfigVariables,
	redisClient *redis.Client,
	kycClient *clients.Client,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.MetricsMiddleware())

	// Expose Prometheus metrics endpoint
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Health check route
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Dynamic routing from config
	middlewareRegistry := middleware.NewRegistry(cfgEnvs, redisClient, kycClient)

	for _, svc := range cfg.Services {
		for _, rt := range svc.Routes {
			handler := discovery.ProxyHandler(svc.BaseURL, rt.ServicePath)

			mws := middlewareRegistry.GetMiddlewares(rt.Middlewares)

			group := r.Group("")
			group.Use(mws...) // custom middleware per route

			switch rt.Method {
			case "GET":
				group.GET(rt.Path, handler)
			case "POST":
				group.POST(rt.Path, handler)
			case "PUT":
				group.PUT(rt.Path, handler)
			case "DELETE":
				group.DELETE(rt.Path, handler)
			}
		}
	}

	return r
}
