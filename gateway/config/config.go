package config

import (
	"os"
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
