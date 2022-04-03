package logger

import (
	"fmt"

	"github.com/seth2810/go-collab/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(cfg *config.LoggerConfig) (*zap.Logger, error) {
	var lvl zapcore.Level

	if err := lvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("failed to set log level: %w", err)
	}

	c := zap.Config{
		Development:      true,
		Encoding:         "console",
		OutputPaths:      []string{cfg.File},
		ErrorOutputPaths: []string{cfg.File},
		Level:            zap.NewAtomicLevelAt(lvl),
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
	}

	return c.Build()
}
