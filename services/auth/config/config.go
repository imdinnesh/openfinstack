package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env        string
	ServerPort string
	DBUrl         string
	RedisUrl      string
	JWTSecret     string
	SMTPHost      string
	SMTPPort      string
	SMTPUser      string
	SMTPPassword  string
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
		ServerPort: getEnv("SERVER_PORT", "8080"),
		DBUrl:         getEnv("DB_URL", "postgres://auth:auth@localhost:5432/fintechdb_auth"),
		RedisUrl:      getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:     getEnv("JWT_SECRET", "supersecretkey"),
		SMTPHost:      getEnv("SMTP_HOST", ""),
		SMTPPort:      getEnv("SMTP_PORT", "587"),
		SMTPUser:      getEnv("SMTP_USER", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
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
