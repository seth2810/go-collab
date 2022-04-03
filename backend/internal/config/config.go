package config

import (
	"fmt"

	"github.com/iamolegga/enviper"
	"github.com/spf13/viper"
)

type LoggerConfig struct {
	Level string `mapstructure:"LOG_LEVEL"`
	File  string `mapstructure:"LOG_FILE"`
}

type HTTPConfig struct {
	Host string `mapstructure:"HTTP_HOST"`
	Port string `mapstructure:"HTTP_PORT" validate:"min=1,max=65535"`
}

type Config struct {
	Logger LoggerConfig `mapstructure:",squash"`
	Server HTTPConfig   `mapstructure:",squash"`
}

func Read() (*Config, error) {
	var cfg Config

	v := enviper.New(viper.New())

	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("while unmarshal config: %w", err)
	}

	if cfg.Logger.File == "" {
		cfg.Logger.File = "/dev/stdout"
	}

	if cfg.Server.Port == "" {
		cfg.Server.Port = "8000"
	}

	return &cfg, nil
}
