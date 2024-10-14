package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DBSource      string        `yaml:"db_source"`
	MigrateSource string        `yaml:"migrate_source"`
	RedisSource   string        `yaml:"redis_source"`
	CacheTTL      time.Duration `yaml:"cache_ttl"`
	HTTPServer    `yaml:"http_server"`
}

type HTTPServer struct {
	Address        string        `yaml:"address" env-default:"0.0.0.0:8080"`
	ReadTimeOut    time.Duration `yaml:"read_timeout"`
	WriteTimeOut   time.Duration `yaml:"write_timeout"`
	IdleTimeout    time.Duration `yaml:"idle_timeout" env-default:"60s"`
	MaxHeaderBytes int           `yaml:"max_header_bytes"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err.Error())
	}

	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err.Error())
	}

	return &cfg
}
