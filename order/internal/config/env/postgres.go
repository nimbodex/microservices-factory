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

func (c *PostgresConfig) Host() string {
	return c.host
}

func (c *PostgresConfig) Port() string {
	return c.port
}

func (c *PostgresConfig) User() string {
	return c.user
}

func (c *PostgresConfig) Password() string {
	return c.password
}

func (c *PostgresConfig) Database() string {
	return c.database
}

func (c *PostgresConfig) SSLMode() string {
	return c.sslMode
}

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
