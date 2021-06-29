package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel          string        `envconfig:"log_level" default:"info"`
	LogLevelFieldName string        `envconfig:"log_level_field_name" default:"severity"`
	ServerAddress     string        `envconfig:"server_address" default:":8080"`
	HTTPTimeout       time.Duration `envconfig:"http_timeout" default:"60s"`
	Hostname          string        `envconfig:"host_name"`
	PrivateKey        string        `envconfig:"private_key"`
	PublicKey         string        `envconfig:"public_key"`
	Ts                string        `envconfig:"ts"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("app", &cfg); err != nil {
		return nil, fmt.Errorf("process config environment variables: %w", err)
	}
	return &cfg, nil
}
