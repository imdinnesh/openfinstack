package main

import (
	"log"

	"github.com/imdinnesh/openfinstack/gateway/config"
	"github.com/imdinnesh/openfinstack/gateway/router"
)

func main() {
	cfg, err := config.LoadConfig("./config/routes.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	r := router.SetupRouter(cfg)
	r.Run(":8000")
}
