package env

import (
	"os"
)

type LoggerConfig struct{}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{}
}

func (c *LoggerConfig) GetLevel() string {
	level := os.Getenv("LOG_LEVEL")
	if level == "" {
		return "info"
	}
	return level
}

func (c *LoggerConfig) GetFormat() string {
	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		return "json"
	}
	return format
}
