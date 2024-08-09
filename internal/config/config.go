package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Env     string `yaml:"env" env-required:"true"`
	Storage `yaml:"storage" env-required:"true"`
	Server  `yaml:"http_server" env-required:"true"`
}

type Server struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"120s"`
}

type Storage struct {
	Host     string `yaml:"POSTGRES_HOST" env-required:"true"`
	Port     int    `yaml:"POSTGRES_PORT" env-required:"true"`
	User     string `yaml:"POSTGRES_USER" env-required:"true"`
	Password string `yaml:"POSTGRES_PASSWORD" env-required:"true"`
	Database string `yaml:"POSTGRES_DB" env-required:"true"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/local.yaml"
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("config file not found: %s", configPath))
	}

	configData, err := os.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var cfg Config
	if err := yaml.Unmarshal(configData, &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
