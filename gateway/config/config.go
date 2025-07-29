package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

type Route struct {
	Path        string   `yaml:"path"`
	Method      string   `yaml:"method"`
	ServicePath string   `yaml:"servicePath"`
	Middlewares []string `yaml:"middlewares"`
}

type Service struct {
	Name    string  `yaml:"name"`
	BaseURL string  `yaml:"baseURL"`
	Routes  []Route `yaml:"routes"`
}

type Config struct {
	Services []Service `yaml:"services"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

type ConfigVariables struct {
	Env          string
	ServerPort   string
	DBUrl        string
	RedisUrl     string
	JWTSecret    string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	KYCBaseURL string
}

func LoadENVS() *ConfigVariables {
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}

	// Load corresponding .env file
	if err := godotenv.Load(".env." + env); err != nil {
		log.Printf("Warning: could not load .env.%s file: %v", env, err)
	}

	cfg := &ConfigVariables{
		Env:        getEnv("ENV", "dev"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisUrl:   getEnv("REDIS_URL", "localhost:6379"),
		JWTSecret:  getEnv("JWT_SECRET", "supersecretkey"),
		DBUrl:         getEnv("DB_URL", "postgres://profile:profile@localhost:5433/fintechdb_kyc"),
		KYCBaseURL: getEnv("KYC_BASE_URL", "http://localhost:8081"),
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
