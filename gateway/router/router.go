package router

import (
	"github.com/gin-gonic/gin"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/gateway/discovery"
	"github.com/imdinnesh/openfinstack/gateway/middleware"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	for _, svc := range cfg.Services {
		for _, rt := range svc.Routes {
			handler := discovery.ProxyHandler(svc.BaseURL, rt.ServicePath)

			// Fetch middlewares for this route
			mws := middleware.GetMiddlewares(rt.Middlewares)

			// Create route group with these middlewares
			group := r.Group("")
			group.Use(mws...)

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
