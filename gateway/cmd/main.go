package main

import (
	"log"

	"github.com/imdinnesh/openfinstack/gateway/clients"
	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/gateway/router"
	"github.com/imdinnesh/openfinstack/packages/redis"
)

func main() {
	cfg, err := config.LoadConfig("./config/routes.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfgEnvs := config.LoadENVS()
	kycClient:=clients.NewClient(cfgEnvs.KYCBaseURL)
	redisClient := redis.NewClient(cfgEnvs.RedisUrl)
	r := router.SetupRouter(cfg,cfgEnvs,redisClient,kycClient)
	r.Run(":8000")
}
