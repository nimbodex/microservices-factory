package env

import (
	"os"
)

type LoggerConfig struct{}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{}
}

func (c *LoggerConfig) Level() string {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		return level
	}
	return "info"
}

func (c *LoggerConfig) AsJSON() bool {
	if format := os.Getenv("LOG_FORMAT"); format != "" {
		return format == "json"
	}
	return false
}
