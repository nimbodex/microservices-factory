package env

import (
	"os"
	"strconv"
)

// LoggerConfig реализация конфигурации логгера из переменных окружения
type LoggerConfig struct {
	level  string
	asJSON bool
}

// NewLoggerConfig создает новую конфигурацию логгера из переменных окружения
func NewLoggerConfig() *LoggerConfig {
	level := getEnvOrDefault("LOGGER_LEVEL", "info")
	asJSONStr := getEnvOrDefault("LOGGER_AS_JSON", "false")
	asJSON, err := strconv.ParseBool(asJSONStr)
	if err != nil {
		asJSON = false
	}

	return &LoggerConfig{
		level:  level,
		asJSON: asJSON,
	}
}

func (c *LoggerConfig) Level() string {
	return c.level
}

func (c *LoggerConfig) AsJSON() bool {
	return c.asJSON
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
