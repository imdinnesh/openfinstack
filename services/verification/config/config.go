package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env        string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
	KYCVerifier string // "manual" "mock", "idfy", or "karza"
	IDFYApiKey string // Used if KYCVerifier is "idfy"
	IDFYBaseURL string // Used if KYCVerifier is "idfy"
	DBUrl     string // Database connection string
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
		SMTPHost:      getEnv("SMTP_HOST", "localhost"),
		SMTPPort:      getEnv("SMTP_PORT", "2525"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
		KYCVerifier:   getEnv("KYC_VERIFIER", "mock"),
		IDFYApiKey:    getEnv("IDFY_API_KEY", ""),
		IDFYBaseURL:   getEnv("IDFY_BASE_URL", "https://api.idfy.com"),
		DBUrl:         getEnv("DB_URL", "postgres://kyc:kyc@localhost:5433/fintechdb_kyc"),
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
