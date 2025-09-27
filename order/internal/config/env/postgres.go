package env

import (
	"fmt"
	"net/url"
)

// PostgresConfig реализация конфигурации PostgreSQL из переменных окружения
type PostgresConfig struct {
	host     string
	port     string
	user     string
	password string
	database string
	sslMode  string
}

// NewPostgresConfig создает новую конфигурацию PostgreSQL из переменных окружения
func NewPostgresConfig() *PostgresConfig {
	return &PostgresConfig{
		host:     getEnvOrDefault("POSTGRES_HOST", "localhost"),
		port:     getEnvOrDefault("POSTGRES_PORT", "5432"),
		user:     getEnvOrDefault("POSTGRES_USER", "postgres"),
		password: getEnvOrDefault("POSTGRES_PASSWORD", ""),
		database: getEnvOrDefault("POSTGRES_DB", "orders"),
		sslMode:  getEnvOrDefault("POSTGRES_SSL_MODE", "disable"),
	}
}

// Host возвращает хост PostgreSQL
func (c *PostgresConfig) Host() string {
	return c.host
}

// Port возвращает порт PostgreSQL
func (c *PostgresConfig) Port() string {
	return c.port
}

// User возвращает имя пользователя для подключения к PostgreSQL
func (c *PostgresConfig) User() string {
	return c.user
}

// Password возвращает пароль для подключения к PostgreSQL
func (c *PostgresConfig) Password() string {
	return c.password
}

// Database возвращает название базы данных
func (c *PostgresConfig) Database() string {
	return c.database
}

// SSLMode возвращает режим SSL
func (c *PostgresConfig) SSLMode() string {
	return c.sslMode
}

// URI возвращает строку подключения к PostgreSQL
func (c *PostgresConfig) URI() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		url.QueryEscape(c.user),
		url.QueryEscape(c.password),
		c.host,
		c.port,
		c.database,
		c.sslMode,
	)
}
