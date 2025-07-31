package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env        string
	DBUrl      string
	RedisUrl   string
	JWTSecret  string
	ServerPort string
	KYCVerifier string // "mock", "idfy", or "karza"
	IDFYApiKey string // Used if KYCVerifier is "idfy"
	IDFYBaseURL string // Used if KYCVerifier is "idfy"
}

func Load() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load corresponding .env file
	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("Warning: could not load .env.%s file: %v", env, err)
	}

	cfg := &Config{
		Env:        getEnv("ENV", "dev"),
		ServerPort: getEnv("SERVER_PORT", "8081"),
		DBUrl:         getEnv("DB_URL", "postgres://profile:profile@localhost:5433/fintechdb_kyc"),
		RedisUrl:      getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:     getEnv("JWT_SECRET", "supersecretkey"),
		KYCVerifier:   getEnv("KYC_VERIFIER", "mock"),
		IDFYApiKey:    getEnv("IDFY_API_KEY", ""),
		IDFYBaseURL:   getEnv("IDFY_BASE_URL", "https://api.idfy.com"),
	}

	fmt.Println("Loaded environment:", cfg.Env)
	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
