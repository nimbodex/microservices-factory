package env

import (
	"net"
)

// OrderHTTPConfig реализация конфигурации HTTP сервера из переменных окружения
type OrderHTTPConfig struct {
	host        string
	port        string
	readTimeout string
}

// NewOrderHTTPConfig создает новую конфигурацию HTTP сервера из переменных окружения
func NewOrderHTTPConfig() *OrderHTTPConfig {
	return &OrderHTTPConfig{
		host:        getEnvOrDefault("HTTP_HOST", "0.0.0.0"),
		port:        getEnvOrDefault("HTTP_PORT", "8080"),
		readTimeout: getEnvOrDefault("HTTP_READ_TIMEOUT", "30s"),
	}
}

// Host возвращает хост HTTP сервера
func (c *OrderHTTPConfig) Host() string {
	return c.host
}

// Port возвращает порт HTTP сервера
func (c *OrderHTTPConfig) Port() string {
	return c.port
}

// ReadTimeout возвращает таймаут чтения HTTP запроса
func (c *OrderHTTPConfig) ReadTimeout() string {
	return c.readTimeout
}

// Address возвращает адрес для прослушивания HTTP сервера
func (c *OrderHTTPConfig) Address() string {
	return net.JoinHostPort(c.host, c.port)
}
