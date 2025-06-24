package main

import (
	Logger "github.com/imdinnesh/openfinstack/packages/logger"
	"net/http"
)

func main() {
	Logger.Log.Info().Msg("Starting Auth Service")
	http.ListenAndServe(":8081", nil)
}
