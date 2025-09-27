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
	asJSON, _ := strconv.ParseBool(asJSONStr)

	return &LoggerConfig{
		level:  level,
		asJSON: asJSON,
	}
}

// Level возвращает уровень логирования
func (c *LoggerConfig) Level() string {
	return c.level
}

// AsJSON возвращает флаг для вывода логов в формате JSON
func (c *LoggerConfig) AsJSON() bool {
	return c.asJSON
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
