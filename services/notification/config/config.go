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
