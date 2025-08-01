package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// KYCVerifierType defines allowed verifier types as an enum-like custom type
type KYCVerifierType string

const (
	KYCVerifierManual KYCVerifierType = "manual"
	KYCVerifierMock   KYCVerifierType = "mock"
	KYCVerifierIDFY   KYCVerifierType = "idfy"
	KYCVerifierKarza  KYCVerifierType = "karza"
)

// IsValid validates if the value is one of the allowed KYCVerifierType values
func (v KYCVerifierType) IsValid() bool {
	switch v {
	case KYCVerifierManual, KYCVerifierMock, KYCVerifierIDFY, KYCVerifierKarza:
		return true
	default:
		return false
	}
}

// Config holds all environment configuration values
type Config struct {
	Env         string
	DBUrl       string
	RedisUrl    string
	JWTSecret   string
	ServerPort  string
	KYCVerifier KYCVerifierType
	IDFYApiKey  string // Only if KYCVerifier is "idfy"
	IDFYBaseURL string // Only if KYCVerifier is "idfy"
}

// Load reads and validates the environment variables, returning the Config
func Load() *Config {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load .env file based on environment
	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("Warning: could not load .env.%s file: %v", env, err)
	}

	// Read and validate KYCVerifier
	verifier := KYCVerifierType(getEnv("KYC_VERIFIER", "mock"))
	if !verifier.IsValid() {
		log.Fatalf("Invalid KYC_VERIFIER value: %s", verifier)
	}

	cfg := &Config{
		Env:         getEnv("ENV", "dev"),
		ServerPort:  getEnv("SERVER_PORT", "8081"),
		DBUrl:       getEnv("DB_URL", "postgres://profile:profile@localhost:5433/fintechdb_kyc"),
		RedisUrl:    getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:   getEnv("JWT_SECRET", "supersecretkey"),
		KYCVerifier: verifier,
		IDFYApiKey:  getEnv("IDFY_API_KEY", ""),
		IDFYBaseURL: getEnv("IDFY_BASE_URL", "https://api.idfy.com"),
	}

	fmt.Println("Loaded environment:", cfg.Env)
	return cfg
}

// getEnv returns the value of the environment variable or a default if not set
func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
